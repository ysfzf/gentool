/*
Copyright © 2022 fzf <ysfzf@hotmail>

*/
package cmd

import (
	"os"

	"github.com/spf13/cobra"
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
