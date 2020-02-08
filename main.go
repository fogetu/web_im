package main

import (
	"github.com/astaxie/beego"
	"github.com/fogetu/web_im/controllers"
)

func main() {
	beego.Router("/", &controllers.MainController{})
	beego.Router("/task/", &controllers.TaskController{}, "get:ListTasks;post:NewTask")
	beego.Router("/upgrade/", &controllers.ChatRoomController{}, "get:Upgrade")
	beego.Router("/task/:id:int", &controllers.TaskController{}, "get:GetTask;put:UpdateTask")
	beego.Run()
}
