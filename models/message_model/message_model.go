package message_model

import (
	"container/list"
	"github.com/fogetu/web_im/models/chat_room_model"
	"sync/atomic"
)

// ROOM消息
// 当有人单独回复某条消息,当前消息所属的UserID用户会被提醒
type MessageModel struct {
	MsgID       uint64
	MsgContent  string // html
	CreateAt    int64
	RoomID      chat_room_model.RoomID
	ParentMsgID uint64 // 是room或者channel下的任意话题，参考slack实现
	UserID      chat_room_model.UserID
}

var RoomMessage = make(map[chat_room_model.RoomID]*list.List)

var (
	IDIncreaseForMsg uint64
)

func GetIncreamID() uint64 {
	atomic.AddUint64(&IDIncreaseForMsg, 1)
	return IDIncreaseForMsg
}
