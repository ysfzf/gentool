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

var rootCmd = &cobra.Command{
	Use:   "gentool",
	Short: "gentool v1.0.2 , Generate a proto file or gorm model file",
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
