package cmd

import (
	"fmt"
	"log"

	"github.com/mynameisfzf/gentool/common"
	"github.com/spf13/cobra"
)

var apiCmd = &cobra.Command{
	Use:   "api",
	Short: "通过数据表创建一个api文件",
	Long: `这个命令可以根据数据表结构创建一个go-zero项目的api文件,栗子:

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
		fmt.Println(" - 未知配置文件")
		return
	}
	var cc common.ProtoConfig
	err := common.LoadConfig(Cfg, &cc)
	if err != nil {
		log.Fatal(err)
	}

	if cc.Schema == "" {
		fmt.Println(" - 未知数据库名 ")
		return
	}

	s, err := cc.GenerateApi()

	if nil != err {
		log.Fatal(err)
	}

	if nil != s {
		writeFile(cc.OutFile, s.String())
	}

	fmt.Println("Done.")

}
