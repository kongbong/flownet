package flownet

import (
	"net"
	"strconv"
)

// Client client network handler
type Client struct {
	RecvChan chan []byte
	DiscChan chan struct{}
	Conn     net.Conn
}

// MakeClient make new client network handler
func MakeClient() *Client {
	client := &Client{}
	client.RecvChan = make(chan []byte)
	client.DiscChan = make(chan struct{})
	return client
}

func (c *Client) onRecv(conn net.Conn, buf []byte) {
	c.RecvChan <- buf
}

func (c *Client) onDisconn(conn net.Conn) {
	c.DiscChan <- struct{}{}
}

// Start client
func (c *Client) Start(port int) error {
	conn, err := net.Dial("tcp", ":"+strconv.Itoa(port))
	if err != nil {
		return err
	}
	c.Conn = conn
	go handleConn(c, conn)
	return nil
}
