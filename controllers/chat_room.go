package controllers

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/fogetu/web_im/models/chat_room_model"
	"github.com/fogetu/web_im/service/chat_events"
	"github.com/fogetu/web_im/service/chat_room_service"
	"github.com/fogetu/web_im/service/websocket"
	"time"
)

type ChatRoomController struct {
	beego.Controller
}

// 进入聊天页面就升级协议
func (chat *ChatRoomController) Upgrade() {
	userID, err := chat.GetInt("user_id")
	if err != nil || userID <= 0 {
		chat.Redirect("/", 302)
		return
	}
	// Upgrade from http request to WebSocket.
	ws, e := websocket.Upgrade(chat.Ctx.ResponseWriter, chat.Ctx.Request, nil, 1024, 1024)
	if e != nil {
		logs.Error("Cannot setup WebSocket connection:", e)
		return
	}
	// Join chat room.
	_, errRet := chat_room_service.UserOnline(chat_room_model.UserID(userID), ws)
	if errRet != nil {
		logs.Error("Cannot setup user online:", errRet)
		return
	}
	defer func() {
		_, err := chat_room_service.UserOffline(chat_room_model.UserID(userID), ws)
		if err != nil {
			logs.Error("Cannot setup user offline:", err)
			return
		}
	}()
	// Message receive loop.
	for {
		_, p, err := ws.ReadMessage()
		fmt.Println(err)
		if err != nil {
			return
		}
		chat_events.Publish <- chat_events.New(chat_events.EVENT_MESSAGE, time.Now().Unix(),
			string(p))
	}
}
