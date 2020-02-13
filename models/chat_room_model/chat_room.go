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

var roomPics = [2]string{"http://106.54.93.177:9091/amodvis/static/image/27/a0/a3/27a0a33aeac4e3b4b8b59a43edb34057.jpeg",
	"http://106.54.93.177:9091/amodvis/static/image/84/de/4c/84de4c5d166dea9a96096ed49d649fa9.jpg"}

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
