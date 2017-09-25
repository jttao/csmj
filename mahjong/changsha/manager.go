package changsha

import "sync"

type RoomManager struct {
	rwm     sync.RWMutex
	roomMap map[int64]*Room
}

func (rm *RoomManager) GetRoomById(roomId int64) *Room {
	rm.rwm.RLock()
	defer rm.rwm.RUnlock()
	r, exist := rm.roomMap[roomId]
	if !exist {
		return nil
	}
	return r
}

func (rm *RoomManager) AddRoom(r *Room) bool {
	rm.rwm.Lock()
	defer rm.rwm.Unlock()
	_, exist := rm.roomMap[r.roomId]
	if exist {
		return false
	}
	rm.roomMap[r.RoomId()] = r
	return true
}

func (rm *RoomManager) RemoveRoom(r *Room) {
	rm.rwm.Lock()
	defer rm.rwm.Unlock()
	delete(rm.roomMap, r.RoomId())
}

func NewRoomManager() *RoomManager {
	rm := &RoomManager{}
	rm.roomMap = make(map[int64]*Room)
	return rm
}
