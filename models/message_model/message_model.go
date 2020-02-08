package message_model

import (
	"container/list"
	"github.com/fogetu/web_im/models/chat_room_model"
)

var MessageData list.List

var RoomMessage = make(map[chat_room_model.RoomID]list.List)
