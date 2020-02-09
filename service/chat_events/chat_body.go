package chat_events

import (
	"github.com/fogetu/web_im/models/chat_room_model"
	"web_im/models/message_model"
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
	message_model.MessageModel
	UserToken string
}

func DecodeToken(token string) chat_room_model.UserID {
	return 0
}
