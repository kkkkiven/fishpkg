package jsonmodel

// ReqGetRoomCardIDs 批量申请房卡号
type ReqGetRoomCardIDs struct {
	RoomID   int64 `json:"roomid"`
	ServerID int64 `json:"serverid"`
	Num      int   `json:"num"`
}

// RespGetRoomCardIDs 批量申请房卡号（响应）
type RespGetRoomCardIDs struct {
	Code int      `json:"code"`
	Msg  string   `json:"msg"`
	IDs  []string `json:"ids"`
}

// ReqReleaseRoomCardIDs 批量释放房卡号
type ReqReleaseRoomCardIDs struct {
	RoomID   int64    `json:"roomid"`
	ServerID int64    `json:"serverid"`
	IDs      []string `json:"ids"`
}
