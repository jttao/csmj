package mahjong

import "time"
import log "github.com/Sirupsen/logrus"

type Heartbeat struct {
	mahjongContext *Mahjong
	done           chan struct{}
}

func (hb *Heartbeat) Start() {
	go func() {
	Loop:
		for {
			select {
			case <-time.After(time.Minute):
				{
					log.Debug("全局处理,定时器")
					now := time.Now().UnixNano() / int64(time.Millisecond)
					for _, pl := range hb.mahjongContext.PlayerManager.Players() {
						elapse := now - pl.LastPing()
						if elapse > int64(time.Minute) {
							log.WithFields(
								log.Fields{
									"玩家id": pl.Id(),
								}).Warn("全局处理,ping超时")
							pl.Session().Close()
						}
					}
				}
			case <-hb.done:
				break Loop
			}
		}
	}()
}

func (hb *Heartbeat) Stop() {
	hb.done <- struct{}{}
}

func NewHeartBeat(m *Mahjong) *Heartbeat {
	return &Heartbeat{
		mahjongContext: m,
	}
}
