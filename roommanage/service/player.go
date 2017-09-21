package service

type Player interface {
	Id() int64
	RoomId() int64
	SetRoomId(roomId int64)
	Ip() string 
	SetIp(ip string) 
}

type player struct {
	id     int64
	roomId int64
	ip     string
}

func (p *player) Id() int64 {
	return p.id
}

func (p *player) RoomId() int64 {
	return p.roomId
}

func (p *player) SetRoomId(roomId int64) {
	p.roomId = roomId
}

func (p *player) Ip() string {
	return p.ip
}

func (p *player) SetIp(ip string) {
	p.ip = ip
}

func NewPlayer(id int64, roomId int64,ip string) Player {
	p := &player{
		id:     id,
		roomId: roomId,
		ip : ip,
	}
	return p
}
