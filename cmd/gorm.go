package cmd

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/mynameisfzf/gentool/common"
	"github.com/spf13/cobra"
	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/gen/internal/generate"
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
		connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", cc.User, cc.Password, cc.Host, cc.Port, cc.Schema)
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

	models := make(map[string]*generate.QueryStructMeta, len(c.Tables))

	for _, tab := range c.Tables {
		// models[tab.Name] = g.GenerateModel(tab.Name)

		// tmp := g.GenerateModelAs(tab.Name, tab.As)
		opts := []gen.FieldOpt{}
		for _, relate := range tab.Relates {
			t := field.BelongsTo

			switch relate.Type {
			case "has_one":
				t = field.HasOne
			case "has_many":
				t = field.HasMany
			case "many_to_many":
				t = field.Many2Many
			}
			model, ok := models[relate.Table]
			if !ok {
				log.Fatal("unknow relate table " + relate.Table)
			}
			opt := gen.FieldRelate(t, common.From(relate.Table).ToCamel(), model, &field.RelateConfig{
				GORMTag: "",
			})
			opts = append(opts, opt)
		}
		if tab.As == "" {
			models[tab.Name] = g.GenerateModel(tab.Name, opts...)
		} else {
			models[tab.Name] = g.GenerateModelAs(tab.Name, tab.As, opts...)
		}
	}
	// models := make([]interface{}, len(c.Tables))

	// for i, table := range c.Tables {
	// 	models[i] = g.GenerateModel(table)
	// }
	if !c.OnlyModel {
		mods := make([]interface{}, len(c.Tables))
		for _, mod := range models {
			mods = append(mods, mod)
		}
		g.ApplyBasic(mods...)

	}
	// shops := g.GenerateModel("shops")
	// fmt.Println(shops.QueryStructName)
	// goods := g.GenerateModel("shop_goods", gen.FieldRelate(field.BelongsTo, "Shop", shops, &field.RelateConfig{}))
	// g.ApplyBasic(shops, goods)
	g.Execute()
	fmt.Println("Done.")

}
