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
	Short: "Generate gorm model file from database",
	Long: `This command generates a Gorm related file, which supports MySQL or Postgres or SQLite or sqlserver. For example:

	gentool gorm --config xx.yaml
 `,
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
		FieldCoverable:    c.FieldCoverable,
		Mode:              gen.WithDefaultQuery | gen.WithQueryInterface,
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
	fmt.Println("Done.")
}
