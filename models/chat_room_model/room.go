package chat_room_model

type OrmUserModel struct {
}

type RoomBaseModel struct {
	RoomID  uint32
	Name    string
	RoomPic string
}

var RoomBaseMap = make(map[int32]*RoomBaseModel)

func init() {
	// mock data
	RoomBaseMap[1] = &RoomBaseModel{1, "room_1", "https://static-upload.local.com/amodvis/static/image/c7/93/54/c793540262a2156d68d10d427a594a02.png"}
	RoomBaseMap[2] = &RoomBaseModel{2, "room_2", "https://static-upload.local.com/amodvis/static/image/92/87/35/928735f507fbc8e411a471210212d028.jpg"}
	RoomBaseMap[3] = &RoomBaseModel{3, "room_3", "https://static-upload.local.com/amodvis/static/image/c7/93/54/c793540262a2156d68d10d427a594a02.png"}
	RoomBaseMap[4] = &RoomBaseModel{4, "room_4", "https://static-upload.local.com/amodvis/static/image/c7/93/54/c793540262a2156d68d10d427a594a02.png"}
	RoomBaseMap[5] = &RoomBaseModel{5, "room_5", "https://static-upload.local.com/amodvis/static/image/c7/93/54/c793540262a2156d68d10d427a594a02.png"}
	RoomBaseMap[6] = &RoomBaseModel{6, "room_6", "https://static-upload.local.com/amodvis/static/image/92/87/35/928735f507fbc8e411a471210212d028.jpg"}
	RoomBaseMap[7] = &RoomBaseModel{7, "room_7", "https://static-upload.local.com/amodvis/static/image/92/87/35/928735f507fbc8e411a471210212d028.jpg"}
	RoomBaseMap[8] = &RoomBaseModel{8, "room_8", "https://static-upload.local.com/amodvis/static/image/c7/93/54/c793540262a2156d68d10d427a594a02.png"}
	RoomBaseMap[9] = &RoomBaseModel{9, "room_9", "https://static-upload.local.com/amodvis/static/image/92/87/35/928735f507fbc8e411a471210212d028.jpg"}
	RoomBaseMap[10] = &RoomBaseModel{10, "room_10", "https://static-upload.local.com/amodvis/static/image/92/87/35/928735f507fbc8e411a471210212d028.jpg"}
}

func (u OrmUserModel) GetByID(roomID int32) *RoomBaseModel {
	if _, ok := RoomBaseMap[roomID]; ok {
		return RoomBaseMap[roomID]
	} else {
		return nil
	}
}
