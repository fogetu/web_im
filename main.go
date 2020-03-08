package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
	_ "github.com/astaxie/beego/migration"
	_ "github.com/astaxie/beego/orm"
	"github.com/fogetu/web_im/controllers"
	_ "github.com/go-sql-driver/mysql"
	"web_im/health_check"
	"web_im/jobs"
)

func init() {
	//跨域设置
	var FilterGateWay = func(ctx *context.Context) {
		origin := ctx.Request.Header.Get("origin")
		logs.Info(origin)
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", origin)
		//允许访问源
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Methods", "POST, GET, PUT, OPTIONS")
		//允许post访问
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Headers", "Access-Control-Allow-Origin,ContentType,Authorization,accept,accept-encoding, authorization, content-type") //header的类型
		ctx.ResponseWriter.Header().Set("Access-Control-Max-Age", "1728000")
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Credentials", "true")
	}
	beego.InsertFilter("*", beego.BeforeRouter, FilterGateWay)
}

func main() {
	health_check.Register()
	jobs.Register()
	beego.Router("/chat_room/get_all_room", &controllers.ChatRoomController{}, "get:GetRoomList;options:GetRoomMessage")
	beego.Router("/chat_room/upgrade", &controllers.ChatRoomController{}, "get:Upgrade")
	beego.Router("/chat_room/create", &controllers.ChatRoomController{}, "get:CreateRoom;post:CreateRoom")
	beego.Router("/chat_room/join", &controllers.ChatRoomController{}, "get:JoinRoom;post:JoinRoom")
	beego.Router("/chat_room/leave", &controllers.ChatRoomController{}, "get:LeaveRoom;post:LeaveRoom")
	beego.Router("/chat_room/room_msg_latest", &controllers.ChatRoomController{}, "get:GetRoomMessage;options:GetRoomMessage")
	beego.Run()
}
