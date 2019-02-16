package flownet

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"flow/common"
	"net"
)

type handler interface {
	onRecv(conn net.Conn, buf []byte)
	onDisconn(conn net.Conn)
}

// SendPacket send packet to conn
func SendPacket(conn net.Conn, data []byte) error {
	if len(data) > MaxPacketSize {
		err := errors.New("writepacket size is too big")
		panic(err)
	}
	if len(data) == 0 {
		return errors.New("writepacket size is zero")
	}

	w := bytes.NewBuffer(make([]byte, 0, 2+len(data)))
	common.WriteData(w, uint16(len(data)))
	common.WriteData(w, data)

	_, err := conn.Write(w.Bytes())
	return err
}

func handleConn(h handler, conn net.Conn) {
	defer conn.Close()

	b := bufio.NewReader(conn)
	readHeader := true
	packetSize := 0
	var headBuf [2]byte
	var packetBuf []byte
	n := 0
	for {
		if readHeader {
			l, err := b.Read(headBuf[n:])
			if err != nil {
				break
			}
			n += l
			if n == 2 {
				packetSize = int(binary.LittleEndian.Uint16(headBuf[:]))
				if packetSize > MaxPacketSize || packetSize == 0 {
					break
				}
				readHeader = false
				n = 0
				packetBuf = make([]byte, packetSize)
			}
		} else {
			l, err := b.Read(packetBuf[n:])
			if err != nil {
				break
			}
			n += l
			if n == packetSize {
				h.onRecv(conn, packetBuf)
				readHeader = true
				n = 0
			}
		}
	}
	h.onDisconn(conn)
}
