package cmd

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	"github.com/mynameisfzf/gentool/common"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "Generate api file from database",
	Long: `This command generates a api file for use in the gozero project. For example:

	gentool api --config xxx.yaml
 `,
	Run: func(cmd *cobra.Command, args []string) {
		generateApi()
	},
}

func init() {
	rootCmd.AddCommand(apiCmd)
}

func generateApi() {
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
	s, err := common.GenerateApi(db, table, cc.IgnoreTables, cc.IgnoreColumns, cc.ServiceName)

	if nil != err {
		log.Fatal(err)
	}

	if nil != s {
		writeFile(cc.OutFile, s.String())
	}

	fmt.Println("Done.")

}
