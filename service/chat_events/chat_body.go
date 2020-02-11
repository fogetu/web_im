package chat_events

import (
	"github.com/fogetu/web_im/models/chat_room_model"
	"web_im/models/message_model"
)

type ChatMsgType uint8

const (
	ChatMsgTypeCommon  ChatMsgType = iota // 普通消息
	ChatMsgTypeJoin                       // 加入ROOM消息
	ChatMsgTypeLeave                      // 退出ROOM消息
	ChatMsgTypeOnline                     // 用户上线
	ChatMsgTypeOffline                    // 用户下线
)

// ROOM消息
// 当有人单独回复某条消息,当前消息所属的UserID用户会被提醒
type MessageBody struct {
	message_model.MessageModel
	UserToken string
}

type MessageJoinRoomBody struct {
	UserID chat_room_model.UserID
	RoomID chat_room_model.RoomID
}

type MessageLeaveRoomBody struct {
	UserID chat_room_model.UserID
	RoomID chat_room_model.RoomID
}

type MessageOnlineBody struct {
	UserID chat_room_model.UserID
}

type MessageOfflineBody struct {
	UserID chat_room_model.UserID
}

// ROOM消息
// 当有人单独回复某条消息,当前消息所属的UserID用户会被提醒
type MessageParentBody struct {
	ChatMstType ChatMsgType
	UserToken   string
	Content     interface{}
}

// ROOM消息
// 当有人单独回复某条消息,当前消息所属的UserID用户会被提醒
type BroadcastMessageBody struct {
	ChatMstType ChatMsgType
	Content     interface{}
}

func DecodeToken(token string) chat_room_model.UserID {
	return 0
}
