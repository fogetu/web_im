package chat_room_model

import (
	"github.com/astaxie/beego/logs"
	"github.com/gorilla/websocket"
	"math/rand"

	//"container/list"
	"sync/atomic"
)

type ChatRoom struct {
	RoomID   RoomID
	RoomName string
	//UserIDList   list.List
	RoomPic        string
	LastMsg        string
	CreateByUserID UserID
	RoomUserList   RoomUserMap
}

type RoomID int32

type UserID uint32

// 用户信息for room
type RoomUserMap map[UserID]UserChatRoomMap

type UserChatRoomMap map[RoomID]UserChatRoom

type ActiveUserInfo struct {
	Conn []*websocket.Conn
}




var (
	IDIncreamForRoom int32
)
var RoomUserList = make(RoomUserMap)
var RoomList = make(map[RoomID]*ChatRoom)
var ActiveUser = make(map[UserID]ActiveUserInfo)

var roomPics = [2]string{"https://static-upload.local.com/amodvis/static/image/c7/93/54/c793540262a2156d68d10d427a594a02.png",
	"https://static-upload.local.com/amodvis/static/image/92/87/35/928735f507fbc8e411a471210212d028.jpg"}

type UserBase struct {
	UserID  uint32
	Name    string
	HeadPic string
}

type UserChatRoom struct {
	RoomID      RoomID
	UnReadCount uint16
	JoinTime    int64
	IsOnline    bool
}

func New(userID UserID, name string) *ChatRoom {
	atomic.AddInt32(&IDIncreamForRoom, 1)
	myRoomID := RoomID(IDIncreamForRoom)
	index := rand.Intn(1)
	RoomList[myRoomID] = &ChatRoom{RoomName: name, RoomID: myRoomID,
		RoomPic: roomPics[index], RoomUserList: RoomUserList, CreateByUserID: userID}
	logs.Info("创建新的ROOM:")
	logs.Info(RoomList)
	return RoomList[myRoomID]
}
