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
	"web_im/models/common_response"
)

type ChatRoomController struct {
	beego.Controller
}

func (chat *ChatRoomController) CreateRoom() {
	userID, err := chat.GetInt("user_id")
	userIDForModel := chat_room_model.UserID(userID)
	if err != nil || userID <= 0 {
		chat.Redirect("/", 302)
		return
	}
	roomName := chat.GetString("room_name")
	if roomName == "" {
		chat.Redirect("/", 302)
		return
	}
	chatRoom, err := chat_room_service.CreateRoom(userIDForModel, roomName)
	if err != nil {
		chat.Data["json"] = common_response.Res{State: false, Msg: "创建失败", Data: false}
		chat.ServeJSON()
	} else {
		_, err = chat_room_service.JoinRoom(userIDForModel, chatRoom.RoomID)
		if err != nil {
			chat.Data["json"] = common_response.Res{State: false, Msg: "创建人加入room失败", Data: false}
			chat.ServeJSON()
		} else {
			chat.Data["json"] = common_response.Res{State: true, Msg: "创建成功", Data: chatRoom}
			chat.ServeJSON()
		}
	}
}

func (chat *ChatRoomController) JoinRoom() {
	roomID, errRoomID := chat.GetInt("room_id")
	userID, err := chat.GetInt("user_id")
	userIDForModel := chat_room_model.UserID(userID)
	roomIDForModel := chat_room_model.RoomID(roomID)
	if errRoomID != nil {
		chat.Redirect("/", 302)
		return
	}
	_, err = chat_room_service.JoinRoom(userIDForModel, roomIDForModel)
	if err != nil {
		chat.Data["json"] = common_response.Res{State: false, Msg: "加入room失败", Data: false}
		chat.ServeJSON()
	} else {
		chat.Data["json"] = common_response.Res{State: true, Msg: "加入room成功", Data: true}
		chat.ServeJSON()
	}
}

func (chat *ChatRoomController) LeaveRoom() {
	roomID, errRoomID := chat.GetInt("room_id")
	userID, err := chat.GetInt("user_id")
	userIDForModel := chat_room_model.UserID(userID)
	roomIDForModel := chat_room_model.RoomID(roomID)
	if errRoomID != nil {
		chat.Redirect("/", 302)
		return
	}
	_, err = chat_room_service.LeaveRoom(userIDForModel, roomIDForModel)
	if err != nil {
		chat.Data["json"] = common_response.Res{State: false, Msg: "退出room失败", Data: false}
		chat.ServeJSON()
	} else {
		chat.Data["json"] = common_response.Res{State: true, Msg: "退出room成功", Data: true}
		chat.ServeJSON()
	}
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
		chat_events.Publish <- chat_events.New(chat_events.EventMessage, time.Now().Unix(),
			string(p))
	}
}
