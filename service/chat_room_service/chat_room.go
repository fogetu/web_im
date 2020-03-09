package chat_room_service

import (
	"errors"
	"github.com/astaxie/beego/logs"
	"github.com/fogetu/web_im/models/chat_room_model"
	ormUer "github.com/fogetu/web_im/models/user_model"
	"github.com/gorilla/websocket"
	"time"
)

func CreateRoom(userID chat_room_model.UserID, roomName string) (*chat_room_model.ChatRoom, error) {
	return chat_room_model.New(userID, roomName), nil
}

// 用户加入ROOM
func JoinRoom(userID chat_room_model.UserID, roomID chat_room_model.RoomID) (bool, error) {
	userIDForModel := userID
	roomIDForModel := roomID
	logs.Info("join romm room_id:", roomIDForModel)
	logs.Info(chat_room_model.RoomList)
	if _, ok := chat_room_model.RoomList[roomIDForModel]; !ok {
		return false, errors.New("room id is not exist")
	}
	OrmUserModel := ormUer.OrmUserModel{}
	userName := OrmUserModel.GetByID(userIDForModel).Name
	// 优先改变刚刚加入ROOM的当前用户状态
	if userMap, ok := chat_room_model.RoomUserList[userIDForModel]; ok {
		if user, ok := userMap[roomIDForModel]; ok {
			user.IsOnline = true
			logs.Info("老用户上线:", userName, "ROOM_ID:", roomID)
		} else {
			userMap[roomIDForModel] = chat_room_model.UserChatRoom{RoomID: chat_room_model.RoomID(roomID),
				JoinTime: time.Now().Unix(), IsOnline: true}
			logs.Info("新用户加入:", userName, "ROOM_ID:", roomID, "FROM:1")
		}
	} else {
		chat_room_model.RoomUserList[userIDForModel] = chat_room_model.UserChatRoomMap{roomIDForModel: chat_room_model.UserChatRoom{RoomID: roomIDForModel, JoinTime: time.Now().Unix(), IsOnline: true}}
		logs.Info("新用户加入:", userName, "ROOM_ID:", roomID, "FROM:2")
	}
	return true, nil
}

// 用户退出ROOM
func LeaveRoom(userID chat_room_model.UserID, roomID chat_room_model.RoomID) (bool, error) {
	if _, ok := chat_room_model.RoomList[roomID]; !ok {
		return false, errors.New("room id is not exist")
	}
	if _, ok := chat_room_model.RoomUserList[userID]; ok {
		delete(chat_room_model.RoomUserList, userID)
	}
	return true, nil
}

// 用户上线
func UserOnline(userID chat_room_model.UserID, ws *websocket.Conn) (bool, error) {

	if _, ok := chat_room_model.ActiveUser[userID]; !ok {
		logs.Info("用户上线,用户未加入任何ROOM:", userID)
		chat_room_model.ActiveUser[userID] = chat_room_model.ActiveUserInfo{Conn: []*websocket.Conn{ws}}
		return true, nil
	}
	connects := chat_room_model.ActiveUser[userID].Conn
	connects = append(connects, ws)
	chat_room_model.ActiveUser[userID] = chat_room_model.ActiveUserInfo{Conn: connects}
	// 他的Join的所有ROOM改变状态
	if _, ok := chat_room_model.RoomUserList[userID]; !ok {
		logs.Info("用户上线,用户未加入任何ROOM:", userID)
		return true, nil
	}
	for key, val := range chat_room_model.RoomUserList[userID] {
		logs.Info("用户上线,更改ROOM用户状态:", "userID:", userID, "roomID:", key)
		val.IsOnline = true
	}
	return true, nil
}

func removeWs(connects []*websocket.Conn, k int) []*websocket.Conn {
	return append(connects[:k], connects[k+1:]...)
}

// 用户下线
func UserOffline(userID chat_room_model.UserID, ws *websocket.Conn) (bool, error) {
	// 他的Join的所有ROOM改变状态
	defer func() {
		connects := chat_room_model.ActiveUser[userID].Conn
		for k, wsCache := range connects {
			if ws == wsCache {
				chat_room_model.ActiveUser[userID] = chat_room_model.ActiveUserInfo{Conn: removeWs(chat_room_model.ActiveUser[userID].Conn, k)}
			}
		}
		err := ws.Close()
		if err != nil {
			logs.Error("close ws error")
		}
	}()
	if _, ok := chat_room_model.RoomUserList[userID]; !ok {
		logs.Info("用户下线,用户未加入任何ROOM:", userID)
		return true, nil
	}
	for key, val := range chat_room_model.RoomUserList[userID] {
		logs.Info("用户下线,更改ROOM用户状态:", "userID:", userID, "roomID:", key)
		val.IsOnline = false
	}
	return true, nil
}
