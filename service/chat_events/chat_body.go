package chat_events

import (
	"github.com/fogetu/web_im/models/chat_room_model"
)

type OnlineBody struct {
	UserID chat_room_model.UserID
}

type OfflineBody struct {
	UserID chat_room_model.UserID
}

// ROOM消息
// 当有人单独回复某条消息,当前消息所属的UserID用户会被提醒
type MessageBody struct {
	MsgID       int64
	MsgContent  string // html
	CreateAt    int64
	RoomID      chat_room_model.RoomID
	ParentMsgID int64 // 是room或者channel下的任意话题，参考slack实现
	UserID      chat_room_model.UserID
}
