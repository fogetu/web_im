package user_model

import (
	"github.com/fogetu/web_im/models/chat_room_model"
	"strconv"
)

var BaseMap = make(map[chat_room_model.UserID]*BaseModel)

type OrmUserModel struct {
}

type BaseModel struct {
	UserID  chat_room_model.UserID
	Name    string
	HeadPic string
}

var MockHeadPic = make([]string, 5, 5)

func init() {
	// mock data
	MockHeadPic[0] = "http://106.54.93.177:9091/amodvis/static/image/fd/83/05/fd8305e6d8cad189dc342aa8ac8aa5a3.jpeg"
	MockHeadPic[1] = "http://106.54.93.177:9091/amodvis/static/image/c6/98/36/c698367fe675370a8d82d513430c6f3e.jpeg"
	MockHeadPic[2] = "http://106.54.93.177:9091/amodvis/static/image/84/de/4c/84de4c5d166dea9a96096ed49d649fa9.jpg"
	MockHeadPic[3] = "http://106.54.93.177:9091/amodvis/static/image/d7/06/9a/d7069a2c5673d3ea42cdb4580c32338a.jpg"
	MockHeadPic[4] = "http://106.54.93.177:9091/amodvis/static/image/64/e2/85/64e2856272980d13e3b60a9f9d46cabf.jpeg"
}

func (u OrmUserModel) GetByID(userID chat_room_model.UserID) *BaseModel {
	index := userID % 5
	var name = `用户` + strconv.FormatUint(uint64(userID), 10)
	userInfo := BaseModel{UserID: userID, Name: name, HeadPic: MockHeadPic[index]}
	return &userInfo
}
