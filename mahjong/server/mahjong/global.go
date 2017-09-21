package mahjong

import log "github.com/Sirupsen/logrus"

type GlobalProcessor struct {
	msgs       chan *Message
	msgHandler MessageHandler
	done       chan struct{}
}

func (gp *GlobalProcessor) Start() {
	go func() {
		log.Info("全局处理器开始")
	Loop:
		for {
			select {
			case msg, ok := <-gp.msgs:
				{
					if !ok {
						break Loop
					}
					err := gp.msgHandler.Handle(msg.Session(), msg.Msg())
					if err != nil {
						log.WithFields(
							log.Fields{
								"error": err,
							}).Warn("全局处理,错误")
					}
				}
			case <-gp.done:
				{
					break Loop
				}
			}
		}
		log.Info("全局处理器结束")
	}()
}

func (gp *GlobalProcessor) Stop() {
	gp.done <- struct{}{}
}

func (gp *GlobalProcessor) Receive(msg *Message) {
	gp.msgs <- msg
}

func NewGlobalProcessor(cacheSize int, msgHandler MessageHandler) *GlobalProcessor {
	gp := &GlobalProcessor{}

	gp.msgs = make(chan *Message, cacheSize)
	gp.done = make(chan struct{})
	gp.msgHandler = msgHandler
	return gp
}
