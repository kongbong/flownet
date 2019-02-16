package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/kongbong/flownet"
)

var cli *flownet.Client

func main() {
	cli = flownet.MakeClient()

	err := cli.Start(port())
	if err != nil {
		panic(err)
	}

	go sendMsg()

	for {
		select {
		case buf := <-cli.RecvChan:
			fmt.Println(string(buf))
		case <-cli.DiscChan:
			log.Println("Disconnected")
			break
		default:
		}
	}
}

func sendMsg() {
	for {
		var msg string
		fmt.Scanln(&msg)
		flownet.SendPacket(cli.Conn, []byte(msg))
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
