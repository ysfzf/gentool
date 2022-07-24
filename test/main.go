package main

import (
	"context"
	"fmt"

	"github.com/mynameisfzf/gentool/dao/query"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func main() {
	connStr := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", "root", "123456", "127.0.0.1", 3306, "fcar")

	db, err := gorm.Open(mysql.Open(connStr))
	if err != nil {
		fmt.Println(err)
		return
	}
	query.SetDefault(db)

	g := query.Q.ShopGoods

	goods, err := g.WithContext(context.Background()).Where(g.ID.Eq(1)).Preload(g.Shops).First()
	if err != nil {
		fmt.Println(err)
		return
	}
	shop, err := g.Shops.WithContext(context.Background()).Model(goods).Find()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Printf("%+v", goods)
	fmt.Println("  ")
	fmt.Println("  ")
	fmt.Println("  ")
	fmt.Printf("%+v", shop)
	g.WithContext(context.Background())
}
