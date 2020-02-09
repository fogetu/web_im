package main

import (
	"github.com/astaxie/beego"
	"github.com/fogetu/web_im/controllers"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/task/", &controllers.TaskController{}, "get:ListTasks;post:NewTask")
	beego.Router("/task/:id:int", &controllers.TaskController{}, "get:GetTask;put:UpdateTask")

	beego.Router("/chat_room/upgrade", &controllers.ChatRoomController{}, "get:Upgrade")
	beego.Router("/chat_room/create", &controllers.ChatRoomController{}, "get:CreateRoom;post:CreateRoom")
	beego.Router("/chat_room/join", &controllers.ChatRoomController{}, "get:JoinRoom;post:JoinRoom")
	beego.Router("/chat_room/leave", &controllers.ChatRoomController{}, "get:LeaveRoom;post:LeaveRoom")

	beego.Run()
}
