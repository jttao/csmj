package api

const (
	RoomManageErrorCode = 30000 + iota
	//玩家已经在房间内
	PlayerAlreadyInRoomErrorCode
	//房间不存在
	RoomNoExistErrorCode
	//加入房间失败
	JoinRoomErrorCode
	//玩家没有在房间内
	PlayerNoInRoomErrorCode
	//房卡不够
	RoomCardNoEnoughErrorCode
	//服务器已满
	ServerIsFullErrorCode
	//账号被封
	AccountIsLockErrorCode
	//服务器即将维护
	ServerWillMaintainedErrorCode
	//定位失败
	LocationErrorCode
	//同 IP 禁止
	JoinRoomIpErrorCode
	//房间人数已经满了
	RoomFullErrorCode
)
