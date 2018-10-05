package main

import (
	"database/sql"
	_ "dfsdbAPI/routers"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	var dataSource="root:lmx1993917@tcp(127.0.0.1:3306)/dfsdb?parseTime=true"
	for true{
		_,err:=sql.Open("mysql",dataSource)
		if err==nil{
			fmt.Println("connected")
			break
		}
	}
	orm.RegisterDataBase("default", "mysql", "root:lmx1993917@tcp(127.0.0.1:3306)/dfsdb")
}

func main() {
	if beego.BConfig.RunMode == "dev" {
		beego.BConfig.WebConfig.DirectoryIndex = true
		beego.BConfig.WebConfig.StaticDir["/swagger"] = "swagger"
	}
	beego.Run()
}

