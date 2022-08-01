package tancy

import (
	"errors"
	"github.com/funny/link"
	"github.com/unsurper/tancy/protocol"
	"sync"
	"sync/atomic"
)

var SessionClosedError = errors.New("Session Closed")
var SessionBlockedError = errors.New("Session Blocked")

var globalSessionId uint64

// 请求上下文
type requestContext struct {
	msgID    uint16
	serialNo uint16
	callback func(answer *protocol.Message)
}

// 终端会话
type Session struct {
	next    uint32
	iccID   uint64
	server  *Server
	session *link.Session

	mux      sync.Mutex
	requests []requestContext

	UserData interface{}
}

// 创建Session
func newSession(server *Server, sess *link.Session) *Session {
	return &Session{
		server:  server,
		session: sess,
	}
}

// 发送消息
func (session *Session) Send(entity protocol.Entity) (uint16, error) {
	message := protocol.Message{
		Body: entity,
		Header: protocol.Header{
			MsgID:       entity.MsgID(),
			IccID:       atomic.LoadUint64(&session.iccID),
			MsgSerialNo: session.nextID(),
		},
	}

	err := session.session.Send(message)
	if err != nil {
		return 0, err
	}
	return message.Header.MsgSerialNo, nil
}

// 获取消息ID
func (session *Session) nextID() uint16 {
	var id uint32
	for {
		id = atomic.LoadUint32(&session.next)
		if id == 0xff {
			if atomic.CompareAndSwapUint32(&session.next, id, 1) {
				id = 1
				break
			}
		} else if atomic.CompareAndSwapUint32(&session.next, id, id+1) {
			id += 1
			break
		}
	}
	return uint16(id)
}

// 回复消息
func (session *Session) Reply(msg *protocol.Message, result protocol.Result) (uint16, error) {
	entity := protocol.T808_0x8001{
		ReplyMsgSerialNo: msg.Header.MsgSerialNo,
		ReplyMsgID:       msg.Header.MsgID,
		Result:           result,
	}
	return session.Send(&entity)
}

// 获取ID
func (session *Session) ID() uint64 {
	return session.session.ID()
}
