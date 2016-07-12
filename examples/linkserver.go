package main

import (
	"flag"
	"os"
	"os/signal"

	"github.com/funny/link"
)

const SECRET = "7#u%AjH63p49"

func main() {
	var addr string

	flag.StringVar(&addr, "addr", ":9090", "echo server address")
	flag.Parse()

	server, err := link.Serve("tcp", addr, link.Async(10, link.Gob()))
	if err != nil {
		panic(err)
	}

	println("server start:", server.Listener().Addr().String())
	go func() {
		for {
			session, err := server.Accept()
			if err != nil {
				break
			}
			go func() {
				addr := session.Conn().RemoteAddr().String()
				println("client", addr, "connected")
				for {
					var msg string
					if err = session.Receive(&msg); err != nil {
						break
					}
					if err = session.Send(msg); err != nil {
						break
					}

				}
				println("client", addr, "closed")
			}()
		}
	}()

	signalChan := make(chan os.Signal)
	signal.Notify(signalChan, os.Interrupt, os.Kill)
	<-signalChan
	server.Stop()

}
