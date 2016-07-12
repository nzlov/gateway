package main

import (
	"flag"
	"fmt"

	"net"
)

const HOSTKEY = "U2FsdGVkX1/rApM8uNHmvPIl9L0t4p6zs4OqLAKu4t0=\n"

func main() {
	var addr string

	flag.StringVar(&addr, "addr", "127.0.0.1:9091", "echo server address")
	flag.Parse()

	tcp4addr, err := net.ResolveTCPAddr("tcp4", addr)
	if err != nil {
		panic(err)
	}
	conn, err := net.DialTCP("tcp4", nil, tcp4addr)
	if err != nil {
		panic(err)
	}
	defer conn.Close()

	_, err = conn.Write([]byte(HOSTKEY))
	if err != nil {
		panic(err)
	}
	quit := make(chan bool)

	go func() {
		msg := make([]byte, 4096)
		for {
			l, err := conn.Read(msg)
			if err != nil {
				fmt.Println("Conn Read Error:", err)
				quit <- true
			}
			fmt.Printf("Recive:%s\n", string(msg[:l]))
		}
	}()

	go func() {
		for {
			var msg string
			if _, err := fmt.Scanf("%s\n", &msg); err != nil {
				break
			}
			if msg == "@quit" {
				quit <- true
				break
			}
			fmt.Printf("Send:%s\n", msg)
			_, err := conn.Write([]byte(msg))
			if err != nil {
				fmt.Println("Conn Write Error:", err)
				quit <- true
				break
			}
		}
	}()
	<-quit
}
