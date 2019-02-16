package main

import (
	"log"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/kongbong/flownet"
)

const tickInterval = 50

var conns []net.Conn

func main() {
	server := flownet.MakeServer()

	server.Start(port())

	conns = make([]net.Conn, 0)
	ticker := time.NewTicker(time.Millisecond * tickInterval)
	lastTickTime := time.Now()

	for {
		select {
		case conn := <-server.ConnChan:
			onConn(conn)
		case packet := <-server.RecvChan:
			onRecv(packet)
		case conn := <-server.DiscChan:
			onDisc(conn)
		case t := <-ticker.C:
			d := t.Sub(lastTickTime)
			onTick(d)
			lastTickTime = t

		}
	}
}

func onConn(conn net.Conn) {
	log.Println("OnConnect", conn)
	conns = append(conns, conn)
}

func onRecv(packet flownet.PacketInfo) {
	log.Println("onRecv", string(packet.Buf))
	broadcast(packet.Buf)

}

func onDisc(conn net.Conn) {
	log.Println("Disconnected", conn)
	for i, c := range conns {
		if c == conn {
			conns = append(conns[:i], conns[i+1:]...)
		}
	}
}

func onTick(d time.Duration) {

}

func broadcast(buf []byte) {
	for _, conn := range conns {
		flownet.SendPacket(conn, buf)
	}
}

func port() int {
	port := os.Getenv("PORT")
	if port != "" {
		p, err := strconv.Atoi(port)
		if err == nil {
			return p
		}
	}
	return 8080
}
