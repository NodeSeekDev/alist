package sftpd

import (
	"golang.org/x/crypto/ssh"
	"net"
)

// Config is the configuration struct for the high level API.
type Config struct {
	// ServerConfig should be initialized properly with
	// e.g. PasswordCallback and AddHostKey
	ssh.ServerConfig
	// HostPort specifies specifies [host]:port to listen on.
	// e.g. ":2022" or "127.0.0.1:2023".
	HostPort string
	// ErrorLogFunc is used to log errors.
	// e.g. log.Println has the right type.
	ErrorLogFunc func(v ...interface{})
	// DebugLogFunc is used to log debug infos.
	// e.g. log.Printf has the right type.
	DebugLogFunc DebugLogger
}

type SftpDriver interface {
	GetConfig() *Config
	GetFileSystem(sc *ssh.ServerConn) (FileSystem, error)
	Close()
}

type SftpServer struct {
	readyChan chan error
	connChan  chan net.Listener
	driver    SftpDriver
}

// NewSftpServer inits a SFTP Server.
func NewSftpServer(driver SftpDriver) *SftpServer {
	return &SftpServer{
		readyChan: make(chan error, 1),
		connChan:  make(chan net.Listener, 1),
		driver:    driver,
	}
}

// RunServer runs the server using the high level API.
func (s *SftpServer) RunServer() error {
	e := runServer(s)
	if e != nil {
		s.LogError("sftpd server failed:", e)
	}
	return e
}

func runServer(server *SftpServer) error {
	listener, e := net.Listen("tcp", server.driver.GetConfig().HostPort)
	server.readyChan <- e
	close(server.readyChan)
	server.connChan <- listener
	close(server.connChan)
	if e != nil {
		return e
	}

	for {
		conn, e := listener.Accept()
		if e != nil {
			return e
		}
		go handleConn(conn, server)
	}
}

func handleConn(conn net.Conn, server *SftpServer) {
	defer func() { _ = conn.Close() }()
	e := doHandleConn(conn, server)
	if e != nil {
		server.LogError("sftpd connection error:", e)
	}
}

func doHandleConn(conn net.Conn, server *SftpServer) error {
	sc, chans, reqs, e := ssh.NewServerConn(conn, &server.driver.GetConfig().ServerConfig)
	if e != nil {
		return e
	}
	defer func() { _ = sc.Close() }()

	// The incoming Request channel must be serviced.
	go printDiscardRequests(server, reqs)

	// Service the incoming Channel channel.
	for newChannel := range chans {
		if newChannel.ChannelType() != "session" {
			_ = newChannel.Reject(ssh.UnknownChannelType, "unknown channel type")
			continue
		}
		channel, requests, err := newChannel.Accept()
		if err != nil {
			return err
		}

		go func(in <-chan *ssh.Request) {
			for req := range in {
				ok := false
				switch {
				case IsSftpRequest(req):
					ok = true
					go func() {
						fs, e := server.driver.GetFileSystem(sc)
						if e == nil {
							var debugf DebugLogger
							if server.driver.GetConfig().DebugLogFunc != nil {
								debugf = server.driver.GetConfig().DebugLogFunc
							} else {
								debugf = func(s string, v ...interface{}) {}
							}
							e = ServeChannel(channel, fs, debugf)
						}
						if e != nil {
							server.LogError("sftpd servechannel failed:", e)
						}
					}()
				}
				_ = req.Reply(ok, nil)
			}
		}(requests)
	}
	return nil
}

func printDiscardRequests(c *SftpServer, in <-chan *ssh.Request) {
	for req := range in {
		c.LogError("sftpd discarding ssh request", req.Type, *req)
		if req.WantReply {
			_ = req.Reply(false, nil)
		}
	}
}

// BlockTillReady will block till the Config is ready to accept connections.
// Returns an error if listening failed. Can be called in a concurrent fashion.
// This is new API - make sure Init is called on the Config before using it.
func (s *SftpServer) BlockTillReady() error {
	err, _ := <-s.readyChan
	return err
}

// Close closes the server assosiated with this config. Can be called in a concurrent
// fashion.
// This is new API - make sure Init is called on the Config before using it.
func (s *SftpServer) Close() error {
	for ch := range s.connChan {
		_ = ch.Close()
	}
	s.driver.Close()
	return nil
}

func (s *SftpServer) LogError(v ...interface{}) {
	if s.driver.GetConfig().ErrorLogFunc != nil {
		s.driver.GetConfig().ErrorLogFunc(v...)
	}
}
