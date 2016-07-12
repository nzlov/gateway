package main

import (
	"flag"
	"fmt"

	"github.com/funny/link"
)

const HOSTKEY = "U2FsdGVkX1/rApM8uNHmvPIl9L0t4p6zs4OqLAKu4t0=\n"

func main() {
	var addr string

	flag.StringVar(&addr, "addr", "127.0.0.1:9091", "echo server address")
	flag.Parse()

	session, err := link.Connect("tcp", addr, link.Async(10, link.Gob()))
	if err != nil {
		panic(err)
	}
	defer session.Close()

	_, err = session.Conn().Write([]byte(HOSTKEY))
	if err != nil {
		panic(err)
	}
	status := make([]byte, 10)
	l, err := session.Conn().Read(status)
	if err != nil {
		panic(err)
	}
	if "200" != string(status[:l]) {
		fmt.Println("连接网关错误！")
		return
	}
	fmt.Println("已经正常连接网关！")

	quit := make(chan bool)

	go func() {
		var msg string
		for {
			if err := session.Receive(&msg); err != nil {
				if err == link.ErrClosed {
					quit <- true
				}
			}
			fmt.Printf("%s\n", msg)
		}
	}()
	go func() {
		for {
			msg := ""
			if _, err := fmt.Scanf("%s", &msg); err != nil {
				break
			}
			if msg == "@quit" {
				quit <- true
				break
			}
			if err = session.Send(msg); err != nil {
				break
			}
		}
	}()

	<-quit
	println("bye")
}
