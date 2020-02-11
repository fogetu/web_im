package main

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/fogetu/web_im/controllers"
)

func init() {
	//跨域设置
	var FilterGateWay = func(ctx *context.Context) {
		ctx.ResponseWriter.Header().Set("Access-Control-Allow-Origin", "http://localhost:4444")
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
	beego.Router("/", &controllers.MainController{})
	beego.Router("/task/", &controllers.TaskController{}, "get:ListTasks;post:NewTask")
	beego.Router("/task/:id:int", &controllers.TaskController{}, "get:GetTask;put:UpdateTask")

	beego.Router("/chat_room/get_all_room", &controllers.ChatRoomController{}, "get:GetRoomList")
	beego.Router("/chat_room/upgrade", &controllers.ChatRoomController{}, "get:Upgrade")
	beego.Router("/chat_room/create", &controllers.ChatRoomController{}, "get:CreateRoom;post:CreateRoom")
	beego.Router("/chat_room/join", &controllers.ChatRoomController{}, "get:JoinRoom;post:JoinRoom")
	beego.Router("/chat_room/leave", &controllers.ChatRoomController{}, "get:LeaveRoom;post:LeaveRoom")
	beego.Run()
}
