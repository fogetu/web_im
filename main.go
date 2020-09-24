package main

import (
	"github.com/astaxie/beego"
	_ "github.com/astaxie/beego/migration"
	_ "github.com/astaxie/beego/orm"
	"github.com/fogetu/web_im/health_check"
	"github.com/fogetu/web_im/jobs"
	"github.com/fogetu/web_im/routers"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	health_check.Register()
	jobs.Register()
	routers.Register()
	beego.Run()
}
