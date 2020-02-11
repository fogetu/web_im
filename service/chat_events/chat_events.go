package chat_events

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/fogetu/web_im/models/chat_room_model"
	"github.com/fogetu/web_im/service/chat_room_service"
	"github.com/gorilla/websocket"
	"time"
	"web_im/models/message_model"
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
	userID := msgData.UserID
	ws := chat_room_model.ActiveUser[userID].Conn
	_, errRet := chat_room_service.UserOnline(userID, ws)
	if errRet != nil {
		logs.Error("Cannot setup user online:", errRet)
	}
	// 广播在线用户
	for _, value := range chat_room_model.ActiveUser {
		ws := value.Conn
		if ws != nil {
			broadcastForOnlineMsg(ws, msgData)
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
	userID := msgData.UserID
	ws := chat_room_model.ActiveUser[userID].Conn
	_, errRet := chat_room_service.UserOnline(userID, ws)
	if errRet != nil {
		logs.Error("Cannot setup user online:", errRet)
	}
	// 广播在线用户
	for _, value := range chat_room_model.ActiveUser {
		ws := value.Conn
		if ws != nil {
			broadcastForOfflineMsg(ws, msgData)
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
		ws := value.Conn
		if ws != nil {
			broadcastForJoinMsg(ws, msgData)
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
		ws := value.Conn
		if ws != nil {
			broadcastForLeaveMsg(ws, msgData)
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
	//msgData.UserID = chat_events.DecodeToken(msgData.UserToken)
	message_model.MessageData.PushBack(msgData)
	// 广播当前ROOM在线用户
	roomUsers := chat_room_model.RoomList[msgData.RoomID].RoomUserList
	for userID, value := range roomUsers {
		isOnline := value[msgData.RoomID].IsOnline
		if isOnline != true {
			continue
		}
		ws := chat_room_model.ActiveUser[userID].Conn
		if ws != nil {
			broadcastForCommonMsg(ws, msgData)
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
	broadcastBody := BroadcastMessageBody{ChatMstType: ChatMsgTypeCommon, Content: msgData}
	b, err := json.Marshal(broadcastBody)
	if err != nil {
		logs.Error(err)
		return
	}
	if ws.WriteMessage(websocket.TextMessage, b) != nil {
		// User disconnected.
		Publish <- New(ChatMsgTypeOffline, time.Now().Unix(), `{"UserID":`+string(msgData.UserID)+`}`)
		logs.Info("ChatMsgType:ChatMsgTypeCommon User disconnected:", msgData.RoomID, "user_id:", msgData.UserID)
	} else {
		logs.Info("ChatMsgType:ChatMsgTypeCommon 广播消息room_id:", msgData.RoomID, "user_id:", msgData.UserID)
	}
}

func New(chatMsgType ChatMsgType, timestamp int64, body string) Event {
	return Event{ChatMsgType: chatMsgType, Timestamp: timestamp, Body: body}
}
