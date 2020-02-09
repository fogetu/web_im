package chat_events

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/fogetu/web_im/models/chat_room_model"
	"github.com/fogetu/web_im/models/message_model"
	"github.com/fogetu/web_im/service/chat_room_service"
	"github.com/gorilla/websocket"
	"time"
)

type ChatEvent uint8

const (
	EventOnline  ChatEvent = iota // 用户上线
	EventOffline                  // 用户下线
	EventMessage                  // 用户发信息
)

type Event struct {
	EventType ChatEvent
	Timestamp int64 // Unix timestamp (secs)
	Body      string
}

var (
	Publish = make(chan Event, 10)
)

func init() {
	go func() {
		for {
			select {
			case event := <-Publish:
				switch event.EventType {
				case EventOnline:
					onlineHandle(event.Body)
				case EventOffline:
					offlineHandle(event.Body)
				case EventMessage:
					messageHandle(event.Body)
				}

			}
		}
	}()
}

func onlineHandle(body string) {
	var msgData OnlineBody
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
}

func offlineHandle(body string) {
	var msgData OfflineBody
	err := json.Unmarshal([]byte(body), &msgData)
	if err != nil {
		logs.Error(err)
		return
	}
	userID := msgData.UserID
	ws := chat_room_model.ActiveUser[userID].Conn
	_, errRet := chat_room_service.UserOffline(userID, ws)
	if errRet != nil {
		logs.Error("Cannot setup user offline:", errRet)
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
			if ws.WriteMessage(websocket.TextMessage, []byte(body)) != nil {
				// User disconnected.
				Publish <- New(EventOffline, time.Now().Unix(), `{"UserID":`+string(userID)+`}`)
			} else {
				logs.Info("广播消息room_id:", msgData.RoomID, "user_id:", userID)
			}
		}
	}
}

func New(eventType ChatEvent, timestamp int64, body string) Event {
	return Event{eventType, timestamp, body}
}
