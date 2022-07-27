package cmd

import (
	"bufio"
	"fmt"
	"log"
	"os"

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
		fmt.Println(" - 未知配置文件")
		return
	}
	var cc common.ProtoConfig
	err := common.LoadConfig(Cfg, &cc)
	if err != nil {
		log.Fatal(err)
	}

	if cc.Schema == "" {
		fmt.Println(" - please input the database schema ")
		return
	}

	s, err := cc.GenerateSchema()

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
