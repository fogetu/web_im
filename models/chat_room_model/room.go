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
	RoomBaseMap[1] = &RoomBaseModel{1, "room_1", "http://106.54.93.177:9091/amodvis/static/image/27/a0/a3/27a0a33aeac4e3b4b8b59a43edb34057.jpeg"}
	RoomBaseMap[2] = &RoomBaseModel{2, "room_2", "http://106.54.93.177:9091/amodvis/static/image/84/de/4c/84de4c5d166dea9a96096ed49d649fa9.jpg"}
	RoomBaseMap[3] = &RoomBaseModel{3, "room_3", "http://106.54.93.177:9091/amodvis/static/image/27/a0/a3/27a0a33aeac4e3b4b8b59a43edb34057.jpeg"}
	RoomBaseMap[4] = &RoomBaseModel{4, "room_4", "http://106.54.93.177:9091/amodvis/static/image/27/a0/a3/27a0a33aeac4e3b4b8b59a43edb34057.jpeg"}
	RoomBaseMap[5] = &RoomBaseModel{5, "room_5", "http://106.54.93.177:9091/amodvis/static/image/27/a0/a3/27a0a33aeac4e3b4b8b59a43edb34057.jpeg"}
	RoomBaseMap[6] = &RoomBaseModel{6, "room_6", "http://106.54.93.177:9091/amodvis/static/image/84/de/4c/84de4c5d166dea9a96096ed49d649fa9.jpg"}
	RoomBaseMap[7] = &RoomBaseModel{7, "room_7", "http://106.54.93.177:9091/amodvis/static/image/84/de/4c/84de4c5d166dea9a96096ed49d649fa9.jpg"}
	RoomBaseMap[8] = &RoomBaseModel{8, "room_8", "http://106.54.93.177:9091/amodvis/static/image/27/a0/a3/27a0a33aeac4e3b4b8b59a43edb34057.jpeg"}
	RoomBaseMap[9] = &RoomBaseModel{9, "room_9", "http://106.54.93.177:9091/amodvis/static/image/84/de/4c/84de4c5d166dea9a96096ed49d649fa9.jpg"}
	RoomBaseMap[10] = &RoomBaseModel{10, "room_10", "http://106.54.93.177:9091/amodvis/static/image/84/de/4c/84de4c5d166dea9a96096ed49d649fa9.jpg"}
}

func (u OrmUserModel) GetByID(roomID int32) *RoomBaseModel {
	if _, ok := RoomBaseMap[roomID]; ok {
		return RoomBaseMap[roomID]
	} else {
		return nil
	}
}
