package models

import (
	//"container/list"
	"sync/atomic"
	"time"
	"github.com/fogetu/web_im"
)

// 用户信息for room
type roomUserMap map[userID]UserChatRoom

type roomID int32
type userID uint32

var (
	RoomUserList     roomUserMap
	IDIncreamForRoom int32
	RoomList         map[roomID]*ChatRoom
)

type ChatRoom struct {
	RoomID   roomID
	RoomName string
	//UserIDList   list.List
	LastMsg      string
	RoomUserList roomUserMap
}

type UserBase struct {
	UserID  uint32
	Name    string
	HeadPic string
}

type UserChatRoom struct {
	RoomID      roomID
	UnReadCount uint16
	JoinTime    int64
	IsOnline    bool
}

func New(name string) *ChatRoom {
	atomic.AddInt32(&IDIncreamForRoom, 1)
	myRoomID := roomID(IDIncreamForRoom)
	RoomList[myRoomID] = &ChatRoom{RoomName: name, RoomID: myRoomID}
	return RoomList[myRoomID]
}

func (r *ChatRoom) Join(userID userID, roomID roomID) {
	if user, ok := r.RoomUserList[userID]; ok {
		user.IsOnline = true
	} else {
		r.RoomUserList[userID] = UserChatRoom{RoomID: roomID, JoinTime: time.Now().Unix(), IsOnline: true}
	}
}

func (r *ChatRoom) Leave(userID userID, roomID roomID) () {
	r.RoomUserList[userID].IsOnline = false
}
