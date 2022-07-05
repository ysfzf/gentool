/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"github.com/mynameisfzf/gentool/core"
	"github.com/spf13/cobra"
	yaml "gopkg.in/yaml.v3"
)

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

// protoCmd represents the proto command
var protoCmd = &cobra.Command{
	Use:   "proto",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		generateProto()
	},
}

func init() {
	rootCmd.AddCommand(protoCmd)
}

func generateProto() {
	if Cfg == "" {
		fmt.Println("未知配置文件")
		return
	}

	cc, err := loadConfig(Cfg)
	if err != nil {
		log.Fatal(err)
	}

	if cc.Schema == "" {
		fmt.Println(" - please input the database schema ")
		return
	}

	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cc.User, cc.Password, cc.Host, cc.Port, cc.Schema)
	db, err := sql.Open(cc.DbType, connStr)
	if err != nil {
		log.Fatal(err)
	}

	defer db.Close()

	table := strings.Join(cc.Tables, ",")
	s, err := core.GenerateSchema(db, table, cc.IgnoreTables, cc.IgnoreColumns, cc.ServiceName, cc.GoPackageName, cc.PackageName)

	if nil != err {
		log.Fatal(err)
	}

	if nil != s {
		//fmt.Println(s)
		writeFile(cc.OutFile, s.String())
	}

}

func loadConfig(path string) (*ProtoConfig, error) {
	var c ProtoConfig
	yamlFile, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	err = yaml.Unmarshal(yamlFile, &c)
	if err != nil {
		return nil, err
	}
	return &c, nil
}

func writeFile(filePath, content string) (bool, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return false, err
	}
	defer file.Close()
	writer := bufio.NewWriter(file) //创建一个writer,带缓存
	writer.WriteString(content)     //写入
	writer.Flush()
	return true, nil
}
