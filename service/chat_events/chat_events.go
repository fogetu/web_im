package chat_events

import (
	"container/list"
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/fogetu/web_im/models/chat_room_model"
	"github.com/fogetu/web_im/service/chat_room_service"
	"github.com/gorilla/websocket"
	"time"
	"web_im/models/message_model"
	"web_im/models/user_model"
)

type Event struct {
	ChatMsgType ChatMsgType
	Timestamp   int64 // Unix timestamp (secs)
	Body        string
}

var (
	Publish = make(chan Event, 10)
)

func init() {
	go func() {
		for {
			select {
			case event := <-Publish:
				switch event.ChatMsgType {
				case ChatMsgTypeOnline:
					onlineHandle(event.Body)
				case ChatMsgTypeOffline:
					logs.Info("-----ChatMsgTypeOfflineChatMsgTypeOffline-----")
					logs.Info(event.Body)
					offlineHandle(event.Body)
				case ChatMsgTypeJoin:
					joinHandle(event.Body)
				case ChatMsgTypeLeave:
					leaveHandle(event.Body)
				case ChatMsgTypeCommon:
					messageHandle(event.Body)
				}
			}
		}
	}()
}

func onlineHandle(body string) {
	var msgData MessageOnlineBody
	err := json.Unmarshal([]byte(body), &msgData)
	if err != nil {
		logs.Error(err)
		return
	}
	// 广播在线用户
	for _, value := range chat_room_model.ActiveUser {
		connects := value.Conn
		if connects != nil {
			for _, ws := range connects {
				broadcastForOnlineMsg(ws, msgData)
			}
		}
	}
}

func offlineHandle(body string) {
	var msgData MessageOfflineBody
	err := json.Unmarshal([]byte(body), &msgData)
	if err != nil {
		logs.Error(err)
		return
	}
	// 广播在线用户
	for _, value := range chat_room_model.ActiveUser {
		connects := value.Conn
		if connects != nil {
			for _, ws := range connects {
				broadcastForOfflineMsg(ws, msgData)
			}
		}
	}
}

func joinHandle(body string) {
	var msgData MessageJoinRoomBody
	err := json.Unmarshal([]byte(body), &msgData)
	if err != nil {
		logs.Error(err)
		return
	}
	_, err = chat_room_service.JoinRoom(msgData.UserID, msgData.RoomID)
	if err != nil {
		logs.Error(err)
		return
	}
	// 广播在线用户
	for _, value := range chat_room_model.ActiveUser {
		connects := value.Conn
		if connects != nil {
			for _, ws := range connects {
				broadcastForJoinMsg(ws, msgData)
			}
		}
	}
}

func leaveHandle(body string) {
	var msgData MessageLeaveRoomBody
	err := json.Unmarshal([]byte(body), &msgData)
	if err != nil {
		logs.Error(err)
		return
	}
	_, err = chat_room_service.LeaveRoom(msgData.UserID, msgData.RoomID)
	if err != nil {
		logs.Error(err)
		return
	}
	// 广播在线用户
	for _, value := range chat_room_model.ActiveUser {
		connects := value.Conn
		if connects != nil {
			for _, ws := range connects {
				broadcastForLeaveMsg(ws, msgData)
			}
		}
	}
}

func messageHandle(body string) {
	var msgData MessageBody
	err := json.Unmarshal([]byte(body), &msgData)
	if err != nil {
		logs.Error(err)
		return
	}
	msgData.MsgID = message_model.GetIncreamID()
	msgData.CreateAt = time.Now().Unix()
	if _, ok := message_model.RoomMessage[msgData.RoomID]; !ok {
		message_model.RoomMessage[msgData.RoomID] = &list.List{}
	}
	message_model.RoomMessage[msgData.RoomID].PushBack(msgData)
	// 广播当前ROOM在线用户
	roomUsers := chat_room_model.RoomList[msgData.RoomID].RoomUserList
	for userID, value := range roomUsers {
		isOnline := value[msgData.RoomID].IsOnline
		if isOnline != true {
			continue
		}
		connects := chat_room_model.ActiveUser[userID].Conn
		if connects != nil {
			for _, ws := range connects {
				broadcastForCommonMsg(ws, msgData)
			}
		}
	}
}

func broadcastForOnlineMsg(ws *websocket.Conn, msgData MessageOnlineBody) {
	broadcastBody := BroadcastMessageBody{ChatMstType: ChatMsgTypeOnline, Content: msgData}
	b, err := json.Marshal(broadcastBody)
	if err != nil {
		logs.Error(err)
		return
	}
	if ws.WriteMessage(websocket.TextMessage, b) != nil {
		// User disconnected.
		logs.Info("ChatMsgType:ChatMsgTypeOnline User disconnected:", "user_id:", msgData.UserID)
	} else {
		logs.Info("ChatMsgType:ChatMsgTypeOnline 广播消息room_id:", "user_id:", msgData.UserID)
	}
}

func broadcastForOfflineMsg(ws *websocket.Conn, msgData MessageOfflineBody) {
	broadcastBody := BroadcastMessageBody{ChatMstType: ChatMsgTypeLeave, Content: msgData}
	b, err := json.Marshal(broadcastBody)
	if err != nil {
		logs.Error(err)
		return
	}
	if ws.WriteMessage(websocket.TextMessage, b) != nil {
		// User disconnected.
		logs.Info("ChatMsgType:ChatMsgTypeLeave User disconnected:", "user_id:", msgData.UserID)
	} else {
		logs.Info("ChatMsgType:ChatMsgTypeLeave 广播消息room_id:", "user_id:", msgData.UserID)
	}
}

func broadcastForJoinMsg(ws *websocket.Conn, msgData MessageJoinRoomBody) {
	broadcastBody := BroadcastMessageBody{ChatMstType: ChatMsgTypeJoin, Content: msgData}
	b, err := json.Marshal(broadcastBody)
	if err != nil {
		logs.Error(err)
		return
	}
	if ws.WriteMessage(websocket.TextMessage, b) != nil {
		// User disconnected.
		logs.Info("ChatMsgType:ChatMsgTypeJoin User disconnected:", msgData.RoomID, "user_id:", msgData.UserID)
	} else {
		logs.Info("ChatMsgType:ChatMsgTypeJoin 广播消息room_id:", msgData.RoomID, "user_id:", msgData.UserID)
	}
}

func broadcastForLeaveMsg(ws *websocket.Conn, msgData MessageLeaveRoomBody) {
	broadcastBody := BroadcastMessageBody{ChatMstType: ChatMsgTypeLeave, Content: msgData}
	b, err := json.Marshal(broadcastBody)
	if err != nil {
		logs.Error(err)
		return
	}
	if ws.WriteMessage(websocket.TextMessage, b) != nil {
		// User disconnected.
		logs.Info("ChatMsgType:ChatMsgTypeLeave User disconnected:", msgData.RoomID, "user_id:", msgData.UserID)
	} else {
		logs.Info("ChatMsgType:ChatMsgTypeLeave 广播消息room_id:", msgData.RoomID, "user_id:", msgData.UserID)
	}
}

func broadcastForCommonMsg(ws *websocket.Conn, msgData MessageBody) {
	type RetMsg struct {
		MessageBody
		UserName string
	}
	OrmUserModel := user_model.OrmUserModel{}
	var retMsg = RetMsg{MessageBody: msgData, UserName: OrmUserModel.GetByID(msgData.UserID).Name}
	broadcastBody := BroadcastMessageBody{ChatMstType: ChatMsgTypeCommon, Content: retMsg}
	b, err := json.Marshal(broadcastBody)
	if err != nil {
		logs.Error(err)
		return
	}
	if ws.WriteMessage(websocket.TextMessage, b) != nil {
		// User disconnected.
		logs.Error("向客户端写信息失败")
		logs.Error(`{"UserID":` + string(msgData.UserID) + `}`)
		Publish <- New(ChatMsgTypeOffline, time.Now().Unix(), `{"UserID":`+string(msgData.UserID)+`}`)
		logs.Info("ChatMsgType:ChatMsgTypeCommon User disconnected:", msgData.RoomID, "user_id:", msgData.UserID)
	} else {
		logs.Info("ChatMsgType:ChatMsgTypeCommon 广播消息room_id:", msgData.RoomID, "user_id:", msgData.UserID)
	}
}

func New(chatMsgType ChatMsgType, timestamp int64, body string) Event {
	return Event{ChatMsgType: chatMsgType, Timestamp: timestamp, Body: body}
}
