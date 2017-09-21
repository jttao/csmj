package changsha_test

import (
	"fmt"
	"math/rand"
	"testing"
	"time"

	"sync"

	. "game/mahjong/changsha"
)

const (
	roomId = 1
)

const (
	initialPlayerId       = 1
	numPlayers            = 4
	waitTime        int64 = int64(1 * time.Second / time.Millisecond)
)

var (
	once       sync.Once
	r          *Room
	roomConfig *RoomConfig = &RoomConfig{
		WaitTime: waitTime,
	}
	customRoomConfig *CustomRoomConfig = &CustomRoomConfig{
		NumPlayers: numPlayers,
	}
)

func setupRoom() {
	once.Do(func() {
		r = NewRoom(roomConfig, customRoomConfig, roomId)
		go func() {
			for {
				select {
				case <-time.After(time.Millisecond * 50):
					r.Tick()
				}
			}
		}()
	})
}

func TestJoinPlayer(t *testing.T) {
	setupRoom()
	rand.Seed(time.Now().UnixNano())
	numRandom := rand.Intn(10) + numPlayers
	fmt.Printf("random join %d players\n", numRandom)
	for i := 0; i < numRandom; i++ {
		p := NewPlayer(initialPlayerId + int64(i))
		flag := r.PlayerJoin(p)
		if i >= numPlayers {
			if flag {
				t.Error("join too much player")
				t.FailNow()
			}
		} else {
			if !flag {
				t.Error("failed to join player")
				t.FailNow()
			}
			defer r.PlayerLeave(p)
		}
	}
}

func TestReconnectPlayer(t *testing.T) {

}

func TestWaitPlayerState(t *testing.T) {
	setupRoom()
	for i := 0; i < numPlayers; i++ {
		p := NewPlayer(initialPlayerId + int64(i))
		flag := r.PlayerJoin(p)

		if !flag {
			t.Error("failed to join player")
			t.FailNow()
		}
	}
	if r.State() != RoomStateWait {
		t.Error("failed enter wait state")
		t.FailNow()
	}
	time.Sleep(time.Duration(1 * waitTime * int64(time.Millisecond)))
	fmt.Print(r.State())
	if r.State() != RoomStatePrepare {
		t.Errorf("failed enter prepare state")
	}
}
