package tancy

import (
	"crypto/rsa"
	"errors"
	"net"
	"strconv"
	"sync"

	"github.com/funny/link"
	log "github.com/sirupsen/logrus"
)

// 服务器选项
type Options struct {
	Keepalive       int64
	AutoMergePacket bool
	CloseHandler    func(*Session)
	PrivateKey      *rsa.PrivateKey
}

// 协议服务器
type Server struct {
	server     *link.Server
	handler    sessionHandler
	timer      *CountdownTimer
	privateKey *rsa.PrivateKey

	mutex    sync.Mutex
	sessions map[uint64]*Session

	closeHandler    func(*Session)
	messageHandlers sync.Map
}

type Handler interface {
	HandleSession(*Session)
}

var _ Handler = HandlerFunc(nil)

type HandlerFunc func(*Session)

func (f HandlerFunc) HandleSession(session *Session) {
	f(session)
}

// 创建服务
func NewServer(options Options) (*Server, error) {
	if options.Keepalive <= 0 {
		options.Keepalive = 60
	}

	if options.PrivateKey != nil && options.PrivateKey.Size() != 128 {
		return nil, errors.New("RSA key must be 1024 bits")
	}

	server := Server{
		closeHandler: options.CloseHandler,
		sessions:     make(map[uint64]*Session),
		privateKey:   options.PrivateKey,
	}
	server.handler.server = &server
	server.handler.autoMergePacket = options.AutoMergePacket
	server.timer = NewCountdownTimer(options.Keepalive, server.handleReadTimeout)
	return &server, nil
}

// 处理读超时
func (server *Server) handleReadTimeout(key string) {
	sessionID, err := strconv.ParseUint(key, 10, 64)
	if err != nil {
		return
	}

	session, ok := server.GetSession(sessionID)
	if !ok {
		return
	}
	session.Close()

	log.WithFields(log.Fields{
		"id":        sessionID,
		"device_id": session.iccID,
	}).Debug("session read timeout")
}

// 获取Session
func (server *Server) GetSession(id uint64) (*Session, bool) {
	server.mutex.Lock()
	defer server.mutex.Unlock()
	session, ok := server.sessions[id]
	if !ok {
		return nil, false
	}
	return session, true
}

func (server *Server) Listener() net.Listener {
	return server.listener
}

func (server *Server) Serve() error {
	for {
		conn, err := Accept(server.listener)
		if err != nil {
			return err
		}

		go func() {
			codec, err := server.protocol.NewCodec(conn)
			if err != nil {
				conn.Close()
				return
			}
			session := server.manager.NewSession(codec, server.sendChanSize)
			server.handler.HandleSession(session)
		}()
	}
}

func (server *Server) Stop() {
	server.listener.Close()
	server.manager.Dispose()
}
