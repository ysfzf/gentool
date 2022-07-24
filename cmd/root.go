/*
Copyright © 2022 fzf <ysfzf@hotmail>

*/
package cmd

import (
	"io/ioutil"
	"os"

	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

var Cfg string

type ProtoConfig struct {
	DbType        string   `yaml:"dbType"`
	Host          string   `yaml:"host"`
	Port          int      `yaml:"port"`
	User          string   `yaml:"user"`
	Password      string   `yaml:"password"`
	Schema        string   `yaml:"schema"`
	Tables        []string `yaml:"tables"`
	ServiceName   string   `yaml:"serviceName"`
	PackageName   string   `yaml:"packageName"`
	GoPackageName string   `yaml:"goPackageName"`
	OutFile       string   `yaml:"outFile"`
	IgnoreTables  []string `yaml:"ignoreTables"`
	IgnoreColumns []string `yaml:"ignoreColumns"`
}

type GenConfig struct {
	DbType            string     `yaml:"dbType"`
	Host              string     `yaml:"host"`
	Port              int        `yaml:"port"`
	User              string     `yaml:"user"`
	Password          string     `yaml:"password"`
	Schema            string     `yaml:"schema"`
	Tables            []GenTable `yaml:"tables"`
	OutPath           string     `yaml:"outPath"`           // specify a directory for output
	OutFile           string     `yaml:"outFile"`           // query code file name, default: gen.go
	OnlyModel         bool       `yaml:"onlyModel"`         // only generate model
	WithUnitTest      bool       `yaml:"withUnitTest"`      // generate unit test for query code
	ModelPkgName      string     `yaml:"modelPkgName"`      // generated model code's package name
	FieldNullable     bool       `yaml:"fieldNullable"`     // generate with pointer when field is nullable
	FieldWithIndexTag bool       `yaml:"fieldWithIndexTag"` // generate field with gorm index tag
	FieldWithTypeTag  bool       `yaml:"fieldWithTypeTag"`  // generate field with gorm column type tag
	FieldCoverable    bool       `yaml:"fieldCoverable"`    //generate pointer when field has default value, to fix problem zero value cannot be assign: https://gorm.io/docs/create.html#Default-Values

}

type GenTable struct {
	Name    string      `yaml:"name"`
	As      string      `yaml:"as"`
	Relates []GenRelate `yaml:"relates"`
}

type GenRelate struct {
	Table  string `yaml:"table"`
	Type   string `yaml:"type"`
	Column string `yaml:"column"`
}

var rootCmd = &cobra.Command{
	Use:   "gentool",
	Short: "gentool v1.0.1 , Generate a proto file or gorm model file",
	Long:  ``,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVarP(&Cfg, "config", "c", "", "config file")

}

func loadConfig[T ProtoConfig | GenConfig](path string, c *T) error {

	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return err
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return err
	}
	return nil
}
