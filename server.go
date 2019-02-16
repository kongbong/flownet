package flownet

import (
	"fmt"
	"net"
	"strconv"
)

// Server network server
type Server struct {
	RecvChan chan PacketInfo
	DiscChan chan net.Conn
	ConnChan chan net.Conn
	conns    map[net.Conn]bool
}

// PacketInfo info of recved packet
type PacketInfo struct {
	Conn net.Conn
	Buf  []byte
}

// MakeServer make new Server
func MakeServer() *Server {
	server := &Server{}
	server.RecvChan = make(chan PacketInfo, 100)
	server.DiscChan = make(chan net.Conn, 100)
	server.ConnChan = make(chan net.Conn, 100)
	server.conns = make(map[net.Conn]bool)
	return server
}

// Start listen
func (s *Server) Start(port int) error {
	server, err := net.Listen("tcp", ":"+strconv.Itoa(port))
	if server == nil {
		return err
	}
	conns := s.clientConns(server)
	go func() {
		for {
			c := <-conns
			s.ConnChan <- c
			go handleConn(s, c)
		}
	}()
	return nil
}

func (s *Server) clientConns(listener net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := listener.Accept()
			if client == nil {
				fmt.Println(err)
				continue
			}
			i++
			fmt.Printf("%d: %v <-> %v\n", i, client.LocalAddr(), client.RemoteAddr())
			s.conns[client] = true
			ch <- client
		}
	}()
	return ch
}

func (s *Server) onRecv(conn net.Conn, buf []byte) {
	s.RecvChan <- PacketInfo{conn, buf}
}

func (s *Server) onDisconn(conn net.Conn) {
	delete(s.conns, conn)
	s.DiscChan <- conn
}
