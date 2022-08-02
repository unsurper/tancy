package tancy

import (
	"github.com/funny/link"
	log "github.com/sirupsen/logrus"
	"strconv"
)

// Session处理
type sessionHandler struct {
	server          *Server
	autoMergePacket bool
}

func (handler sessionHandler) HandleSession(sess *link.Session) {
	log.WithFields(log.Fields{
		"id": sess.ID(),
	}).Debug("[tancy-flow] new session created")

	// 创建Session
	session := newSession(handler.server, sess)
	handler.server.mutex.Lock()
	handler.server.sessions[sess.ID()] = session
	handler.server.mutex.Unlock()
	handler.server.timer.Update(strconv.FormatUint(session.ID(), 10))
	sess.AddCloseCallback(nil, nil, func() {
		handler.server.handleClose(session)
	})

	for {
		// 接收消息
		_, err := sess.Receive()
		if err != nil {
			sess.Close()
			break
		}

	}
}
