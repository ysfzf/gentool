package common

import (
	"bytes"
	"database/sql"
	"fmt"
	"regexp"
	"strings"

	"github.com/chuckpreslar/inflect"
	"github.com/serenize/snaker"
)

type TabField struct {
	Typ     string `json:"typ"`
	Name    string `json:"name"`
	Comment string `json:"comment"`
}

type Tab struct {
	Name    string     `json:"name"`
	Comment string     `json:"comment"`
	Fields  []TabField `json:"fields"`
}

type SchemaApi struct {
	Syntax      string         `json:"syntax"`
	ServiceName string         `json:"serivce"`
	Tables      []Tab          `json:"tables"`
	Enums       EnumCollection `json:"enums,omitempty"`
}

func GenerateApi(db *sql.DB, table string, ignoreTables, ignoreColumns []string, serviceName string) (*SchemaApi, error) {
	s := &SchemaApi{}
	except = append(ignoreColumns, "id")
	dbs, err := dbSchema(db)
	if nil != err {
		return nil, err
	}

	s.Syntax = "v1"
	s.ServiceName = serviceName

	cols, err := dbColumns(db, dbs, table)
	if nil != err {
		return nil, err
	}

	typMap := map[string]*Tab{}
	ignoreMap := map[string]bool{}
	for _, ig := range ignoreTables {
		ignoreMap[ig] = true
	}

	for _, col := range cols {
		if _, ok := ignoreMap[col.TableName]; ok {
			continue
		}

		name := snaker.SnakeToCamel(col.TableName)
		name = inflect.Singularize(name)

		msg, ok := typMap[name]
		if !ok {
			typMap[name] = &Tab{Name: name, Comment: col.TableComment}
			msg = typMap[name]
		}

		err := parseApiColumn(s, msg, col)
		if nil != err {
			return nil, err
		}
		// s.Types = append(s.Types, *msg)
	}

	for _, t := range typMap {
		s.Tables = append(s.Tables, *t)
	}

	return s, nil
}

func parseApiColumn(s *SchemaApi, msg *Tab, col Column) error {
	typ := strings.ToLower(col.DataType)
	var fieldType string

	switch typ {
	case "char", "varchar", "text", "longtext", "mediumtext", "tinytext":
		fieldType = "string"
	case "enum", "set":
		// Parse c.ColumnType to get the enum list
		enumList := regexp.MustCompile(`[enum|set]\((.+?)\)`).FindStringSubmatch(col.ColumnType)
		enums := strings.FieldsFunc(enumList[1], func(c rune) bool {
			cs := string(c)
			return "," == cs || "'" == cs
		})

		enumName := inflect.Singularize(snaker.SnakeToCamel(col.TableName)) + snaker.SnakeToCamel(col.ColumnName)
		enum, err := newEnumFromStrings(enumName, col.ColumnComment, enums)
		if nil != err {
			return err
		}

		s.Enums = append(s.Enums, enum)

		fieldType = enumName
	case "blob", "mediumblob", "longblob", "varbinary", "binary":
		fieldType = "bytes"
	case "date", "time", "datetime", "timestamp":

		fieldType = "string"
	case "bool":
		fieldType = "bool"
	case "tinyint", "smallint", "int", "mediumint", "bigint":
		fieldType = "int64"
	case "float", "decimal", "double":
		fieldType = "float64"
	}

	if "" == fieldType {
		return fmt.Errorf("no compatible protobuf type found for `%s`. column: `%s`.`%s`", col.DataType, col.TableName, col.ColumnName)
	}

	field := TabField{
		Typ:     fieldType,
		Name:    col.ColumnName,
		Comment: col.ColumnComment,
	}

	msg.Fields = append(msg.Fields, field)

	return nil
}

func (s *SchemaApi) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString(fmt.Sprintf("syntax = \"%s\";\n", s.Syntax))
	buf.WriteString("\n")
	buf.WriteString("info( \n")
	buf.WriteString("\t title: \" " + s.ServiceName + " \"\n")
	buf.WriteString("\t desc: \"API文件\"\n")
	buf.WriteString("\t author: \"xxx\"\n")
	buf.WriteString("\t email: \"xxx@yyy.com\"\n")
	buf.WriteString("\t version: \"v1\"\n")
	buf.WriteString(")\n\n")

	buf.WriteString("// ------------------------------------ \n")
	buf.WriteString("// Types\n")
	buf.WriteString("// ------------------------------------ \n\n")
	for _, tab := range s.Tables {
		buf.WriteString("//--------------------------------" + tab.Comment + "--------------------------------")
		buf.WriteString("\n")
		tab.genDefault(buf)
		tab.genGetAll(buf)
	}

	return buf.String()
}

func (tab Tab) genDefault(buf *bytes.Buffer) {
	buf.WriteString("type " + tab.Name + " {\n")
	for _, field := range tab.Fields {
		name := From(field.Name).ToCamel()
		comment := ""
		tag := fmt.Sprintf("`gorm:\"column:%s\" json:\"%s\"`", field.Name, field.Name)
		if field.Comment != "" {
			comment = "// " + field.Comment
		}
		buf.WriteString(fmt.Sprintf("   %s  %s %s  %s \n", name, field.Typ, tag, comment))
	}

	buf.WriteString("}\n")
}

func (tab Tab) genGetAll(buf *bytes.Buffer) {
	buf.WriteString("type Get" + tab.Name + "Request {\n")
	for _, field := range tab.Fields {
		if !isInSlice(except, field.Name) {
			name := From(field.Name).ToCamel()
			comment := ""
			tag := fmt.Sprintf("`form:\"column:%s,optional\"`", field.Name)
			if field.Comment != "" {
				comment = "// " + field.Comment
			}
			buf.WriteString(fmt.Sprintf("   %s  %s %s  %s \n", name, field.Typ, tag, comment))
		}

	}

	buf.WriteString("}\n")
}
