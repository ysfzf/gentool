package cmd

import (
	"bufio"
	"database/sql"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/mynameisfzf/gentool/common"
	"github.com/spf13/cobra"
)

var protoCmd = &cobra.Command{
	Use:   "proto",
	Short: "Generate protobuf file from database",
	Long: `This command generates a proto file for use in the gozero project. For example:

	gentool proto --config xxx.yaml
 `,
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
	s, err := common.GenerateSchema(db, table, cc.IgnoreTables, cc.IgnoreColumns, cc.ServiceName, cc.GoPackageName, cc.PackageName)

	if nil != err {
		log.Fatal(err)
	}

	if nil != s {

		writeFile(cc.OutFile, s.String())
	}

	fmt.Println("Done.")

}

func writeFile(filePath, content string) (bool, error) {
	file, err := os.OpenFile(filePath, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		return false, err
	}
	defer file.Close()
	writer := bufio.NewWriter(file)
	writer.WriteString(content)
	writer.Flush()
	return true, nil
}
