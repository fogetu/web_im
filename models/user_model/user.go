package user_model

var BaseMap = make(map[int32]*BaseModel)

type OrmUserModel struct {
}

type BaseModel struct {
	UserID  uint32
	Name    string
	HeadPic string
}

func init() {
	// mock data
	BaseMap[1] = &BaseModel{1, "user_1", "https://static-upload.local.com/amodvis/static/image/14/9c/ff/149cffa10522fed6855612a647924663.jpeg"}
	BaseMap[2] = &BaseModel{2, "user_2", "https://static-upload.local.com/amodvis/static/image/bc/5a/1c/bc5a1c00e7bbda9074e6626f12e8c0ac.jpeg"}
	BaseMap[3] = &BaseModel{3, "user_3", "https://static-upload.local.com/amodvis/static/image/54/3e/2f/543e2f63c6ed1a043757c2e3c2d975c5.jpeg"}
	BaseMap[4] = &BaseModel{4, "user_4", "https://static-upload.local.com/amodvis/static/image/14/9c/ff/149cffa10522fed6855612a647924663.jpeg"}
	BaseMap[5] = &BaseModel{5, "user_5", "https://static-upload.local.com/amodvis/static/image/14/9c/ff/149cffa10522fed6855612a647924663.jpeg"}
	BaseMap[6] = &BaseModel{6, "user_6", "https://static-upload.local.com/amodvis/static/image/14/9c/ff/149cffa10522fed6855612a647924663.jpeg"}
	BaseMap[7] = &BaseModel{7, "user_7", "https://static-upload.local.com/amodvis/static/image/bc/5a/1c/bc5a1c00e7bbda9074e6626f12e8c0ac.jpeg"}
	BaseMap[8] = &BaseModel{8, "user_8", "https://static-upload.local.com/amodvis/static/image/54/3e/2f/543e2f63c6ed1a043757c2e3c2d975c5.jpeg"}
	BaseMap[9] = &BaseModel{9, "user_9", "https://static-upload.local.com/amodvis/static/image/bc/5a/1c/bc5a1c00e7bbda9074e6626f12e8c0ac.jpeg"}
	BaseMap[10] = &BaseModel{10, "user_10", "https://static-upload.local.com/amodvis/static/image/bc/5a/1c/bc5a1c00e7bbda9074e6626f12e8c0ac.jpeg"}
}

func (u OrmUserModel) GetByID(userID int32) *BaseModel {
	if _, ok := BaseMap[userID]; ok {
		return BaseMap[userID]
	} else {
		return nil
	}
}
