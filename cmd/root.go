/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

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
	DbType            string   `yaml:"dbType"`
	Host              string   `yaml:"host"`
	Port              int      `yaml:"port"`
	User              string   `yaml:"user"`
	Password          string   `yaml:"password"`
	Schema            string   `yaml:"schema"`
	Tables            []string `yaml:"tables"`
	OutPath           string   `yaml:"outPath"`           // specify a directory for output
	OutFile           string   `yaml:"outFile"`           // query code file name, default: gen.go
	OnlyModel         bool     `yaml:"onlyModel"`         // only generate model
	WithUnitTest      bool     `yaml:"withUnitTest"`      // generate unit test for query code
	ModelPkgName      string   `yaml:"modelPkgName"`      // generated model code's package name
	FieldNullable     bool     `yaml:"fieldNullable"`     // generate with pointer when field is nullable
	FieldWithIndexTag bool     `yaml:"fieldWithIndexTag"` // generate field with gorm index tag
	FieldWithTypeTag  bool     `yaml:"fieldWithTypeTag"`  // generate field with gorm column type tag
}

// type Config interface{

// 	ProtoConfig | GenConfig

// }

var rootCmd = &cobra.Command{
	Use:   "gentool",
	Short: "A brief description of your application",
	Long: `A longer description that spans multiple lines and likely contains
examples and usage of using your application. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {

	rootCmd.PersistentFlags().StringVar(&Cfg, "config", "", "config file (default is $HOME/.gentool.yaml)")

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
