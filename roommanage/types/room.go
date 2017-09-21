package types

type Room struct {
	Id       int64
	RoomType RoomType
	OwnerId  int64
	Players  []int64
	Config   string
}

type RoomType int

const (
	RoomTypeChangSha RoomType = iota
)
