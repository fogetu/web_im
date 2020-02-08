package chat_events

import (
	"encoding/json"
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"time"
	"github.com/fogetu/web_im/models/chat_room_model"
	"github.com/fogetu/web_im/models/message_model"
	"github.com/fogetu/web_im/service/chat_room_service"
)

type ChatEvent uint8

const (
	EVENT_ONLINE  ChatEvent = iota // 用户上线
	EVENT_OFFLINE                  // 用户下线
	EVENT_MESSAGE                  // 用户发信息
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
				case EVENT_ONLINE:
					onlineHandle(event.Body)
				case EVENT_OFFLINE:
					offlineHandle(event.Body)
				case EVENT_MESSAGE:
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
				Publish <- New(EVENT_OFFLINE, time.Now().Unix(), `{"UserID":`+string(userID)+`}`)
			}
		}
	}
}

func New(eventType ChatEvent, timestamp int64, body string) Event {
	return Event{eventType, timestamp, body}
}
