package controllers

import (
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/fogetu/web_im/models/chat_room_model"
	"github.com/fogetu/web_im/service/chat_events"
	"github.com/fogetu/web_im/service/chat_room_service"
	"github.com/fogetu/web_im/service/websocket"
	"time"
	"web_im/models/common_response"
	"web_im/models/message_model"
	"web_im/models/user_model"
)

type ChatRoomController struct {
	beego.Controller
}

type allRoomListItem struct {
	RoomID   chat_room_model.RoomID
	RoomName string
	//UserIDList   list.List
	RoomPic        string
	LastMsg        string
	CreateByUserID chat_room_model.UserID
}

func (chat *ChatRoomController) GetRoomList() {
	var allRoomList []allRoomListItem
	for _, v := range chat_room_model.RoomList {
		allRoomList = append(allRoomList, allRoomListItem{RoomID: v.RoomID, RoomName: v.RoomName, RoomPic: v.RoomPic, LastMsg: v.LastMsg, CreateByUserID: v.CreateByUserID})
	}
	chat.Data["json"] = common_response.Res{State: true, Msg: "ok", Data: allRoomList}
	chat.ServeJSON()
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

func (chat *ChatRoomController) GetRoomMessage() {
	roomID, err := chat.GetInt("room_id")
	if err != nil {
		chat.Redirect("/", 302)
		return
	}
	_, err = chat.GetInt("page")
	if err != nil {
		chat.Redirect("/", 302)
		return
	}
	_, err = chat.GetInt("page_size")
	if err != nil {
		chat.Redirect("/", 302)
		return
	}
	if data, ok := message_model.RoomMessage[chat_room_model.RoomID(roomID)]; ok {
		type RetMsg struct {
			chat_events.MessageBody
			UserName string
		}
		var allMsg []RetMsg
		for item := data.Front(); nil != item; item = item.Next() {
			temUser := item.Value.(chat_events.MessageBody)
			OrmUserModel := user_model.OrmUserModel{}
			var retMsg = RetMsg{MessageBody: temUser, UserName: OrmUserModel.GetByID(temUser.UserID).Name}
			allMsg = append(allMsg, retMsg)
		}
		chat.Data["json"] = common_response.Res{State: true, Msg: "success", Data: allMsg}
		chat.ServeJSON()
	} else {
		var allMsg []chat_events.MessageBody
		chat.Data["json"] = common_response.Res{State: true, Msg: "无数据", Data: allMsg}
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
		_, body, err := ws.ReadMessage()
		fmt.Println(err)
		if err != nil {
			logs.Info("receive------skip:", err)
			return
		}
		logs.Info("receive------message")
		logs.Info(string(body))
		var msgParentDataOrigin chat_events.MessageParentBody
		err = json.Unmarshal([]byte(body), &msgParentDataOrigin)
		if err != nil {
			logs.Error(err)
			continue
		}
		switch msgParentDataOrigin.ChatMstType {
		case chat_events.ChatMsgTypeCommon:
			var msgParentData struct {
				chat_events.MessageParentBody
				Content chat_events.MessageBody
			}
			err = json.Unmarshal([]byte(body), &msgParentData)
			if err != nil {
				logs.Error("message------1")
				break
			}
			data, err := json.Marshal(msgParentData.Content)
			if err != nil {
				logs.Error("message------2")
				break
			}
			if _, ok := chat_room_model.RoomList[msgParentData.Content.RoomID]; !ok {
				logs.Error("message------3")
				break
			}
			chat_events.Publish <- chat_events.New(msgParentData.ChatMstType, time.Now().Unix(),
				string(data))
		case chat_events.ChatMsgTypeJoin:
			var msgParentData struct {
				chat_events.MessageParentBody
				Content chat_events.MessageJoinRoomBody
			}
			err = json.Unmarshal([]byte(body), &msgParentData)
			if err != nil {
				logs.Error("message------4")
				break
			}
			data, err := json.Marshal(msgParentData.Content)
			if err != nil {
				logs.Error("message------5")
				break
			}
			chat_events.Publish <- chat_events.New(msgParentData.ChatMstType, time.Now().Unix(),
				string(data))
		case chat_events.ChatMsgTypeLeave:
			var msgParentData struct {
				chat_events.MessageParentBody
				Content chat_events.MessageLeaveRoomBody
			}
			err = json.Unmarshal([]byte(body), &msgParentData)
			if err != nil {
				logs.Error("message------6")
				break
			}
			data, err := json.Marshal(msgParentData.Content)
			if err != nil {
				logs.Error("message------7")
				break
			}
			logs.Info("-----------user leave event-----------")
			logs.Info(string(data))
			chat_events.Publish <- chat_events.New(msgParentData.ChatMstType, time.Now().Unix(),
				string(data))
		case chat_events.ChatMsgTypeOnline:
			var msgParentData struct {
				chat_events.MessageParentBody
				Content chat_events.MessageOnlineBody
			}
			err = json.Unmarshal([]byte(body), &msgParentData)
			if err != nil {
				logs.Error("message------8")
				break
			}
			data, err := json.Marshal(msgParentData.Content)
			if err != nil {
				logs.Error("message------9")
				break
			}
			chat_events.Publish <- chat_events.New(msgParentData.ChatMstType, time.Now().Unix(),
				string(data))
		case chat_events.ChatMsgTypeOffline:
			var msgParentData struct {
				chat_events.MessageParentBody
				Content chat_events.MessageOfflineBody
			}
			err = json.Unmarshal([]byte(body), &msgParentData)
			if err != nil {
				logs.Error("message------10")
				break
			}
			data, err := json.Marshal(msgParentData.Content)
			if err != nil {
				logs.Error("message------11")
				break
			}
			chat_events.Publish <- chat_events.New(msgParentData.ChatMstType, time.Now().Unix(),
				string(data))
		default:
			break
		}
	}
}
