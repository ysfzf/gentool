package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mynameisfzf/gentool/common"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/gorm"
)

var gormCmd = &cobra.Command{
	Use:   "gorm",
	Short: "Generate gorm related file from database",
	Long: `This command generates a Gorm related file, which supports MySQL or Postgres or SQLite or sqlserver. For example:

	gentool gorm --config xx.yaml
 `,
	Run: func(cmd *cobra.Command, args []string) {
		var cc common.GenConfig

		err := common.LoadConfig(Cfg, &cc)
		if err != nil {
			log.Fatal(err)
		}
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cc.User, cc.Password, cc.Host, cc.Port, cc.Schema)
		db, err := sql.Open(cc.DbType, connStr)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()
		generateGorm(db, &cc)
	},
}

func init() {
	rootCmd.AddCommand(gormCmd)
}

func generateGorm(db *sql.DB, c *common.GenConfig) {
	g := gen.NewGenerator(gen.Config{
		OutPath:           c.OutPath,
		OutFile:           c.OutFile,
		ModelPkgPath:      c.ModelPkgName,
		WithUnitTest:      c.WithUnitTest,
		FieldNullable:     c.FieldNullable,
		FieldWithIndexTag: c.FieldWithIndexTag,
		FieldWithTypeTag:  c.FieldWithTypeTag,
		FieldCoverable:    c.FieldCoverable,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
	})

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}))
	if nil != err {
		log.Fatal(err)
	}
	g.UseDB(gdb)

	dataMap := map[string]func(detailType string) (dataType string){
		"int":     func(detailType string) (dataType string) { return "int64" },
		"tinyint": func(detailType string) (dataType string) { return "int8" },
		"date":    func(detailType string) (dataType string) { return "string" },
		"time":    func(detailType string) (dataType string) { return "string" },
	}

	g.WithDataTypeMap(dataMap)
	models := make([]interface{}, len(c.Tables))

	for index, tab := range c.Tables {

		opts := []gen.FieldOpt{
			gen.FieldType("id", "int64"),
			gen.FieldType("deleted_at", "gorm.DeletedAt"),
		}
		for _, relate := range tab.Relates {

			if relate.Column == "" {
				log.Fatal("unkonw relate column")
			}

			if relate.Table == "" {
				log.Fatal("unkonw relate table")
			}

			t := field.BelongsTo

			switch relate.Type {
			case "has_one":
				t = field.HasOne
			case "has_many":
				t = field.HasMany
			case "many_to_many":
				t = field.Many2Many
			}

			tmpModel := g.GenerateModel(relate.Table)
			opt := gen.FieldRelate(t, common.From(relate.Table).ToCamel(), tmpModel, &field.RelateConfig{
				GORMTag: "foreignKey:" + relate.Column,
			})
			opts = append(opts, opt)
		}
		if tab.As == "" {
			models[index] = g.GenerateModel(tab.Name, opts...)
		} else {
			models[index] = g.GenerateModelAs(tab.Name, tab.As, opts...)
		}
	}

	if !c.OnlyModel {
		g.ApplyBasic(models...)

	}

	g.Execute()
	fmt.Println("Done.")

}
