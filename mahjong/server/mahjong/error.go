package mahjong

type ErrorCode int32

const (
	//登陆超时
	ErrorCodeAuthTimeout ErrorCode = iota + 1
	//ping超时
	PingTimeout
	//房间不存在
	ErrorRoomNotFound 
)
