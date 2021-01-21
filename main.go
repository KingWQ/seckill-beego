package main

import (
	"github.com/astaxie/beego/orm"
	_ "shop-seckill/routers"

	"github.com/astaxie/beego"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	defaultdb := beego.AppConfig.String("defaultdb")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	orm.RegisterDataBase("default", "mysql", defaultdb, 30, 30)

	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}
