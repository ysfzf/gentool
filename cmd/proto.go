/*
Copyright © 2022 NAME HERE <EMAIL ADDRESS>

*/
package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mynameisfzf/gentool/core"
	"github.com/spf13/cobra"
)

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
	var cc ProtoConfig
	err := loadConfig(Cfg, &cc)
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
