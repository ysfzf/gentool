package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// gormCmd represents the gorm command
var gormCmd = &cobra.Command{
	Use:   "gorm",
	Short: "A brief description of your command",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		var cc GenConfig

		err := loadConfig(Cfg, &cc)
		if err != nil {
			log.Fatal(err)
		}
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s", cc.User, cc.Password, cc.Host, cc.Port, cc.Schema)
		db, err := sql.Open(cc.DbType, connStr)
		if err != nil {
			log.Fatal(err)
		}

		defer db.Close()
		generateGorm(db, &cc)
	},
}

func init() {
	rootCmd.AddCommand(gormCmd)
}

func generateGorm(db *sql.DB, c *GenConfig) {
	g := gen.NewGenerator(gen.Config{
		OutPath:           c.OutPath,
		OutFile:           c.OutFile,
		ModelPkgPath:      c.ModelPkgName,
		WithUnitTest:      c.WithUnitTest,
		FieldNullable:     c.FieldNullable,
		FieldWithIndexTag: c.FieldWithIndexTag,
		FieldWithTypeTag:  c.FieldWithTypeTag,
	})

	gdb, err := gorm.Open(mysql.New(mysql.Config{
		Conn: db,
	}))
	if nil != err {
		log.Fatal(err)
	}
	g.UseDB(gdb)
	models := make([]interface{}, len(c.Tables))
	for i, table := range c.Tables {
		models[i] = g.GenerateModel(table)
	}
	if !c.OnlyModel {
		g.ApplyBasic(models...)
	}

	g.Execute()
}
