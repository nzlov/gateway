package main

import (
	"fmt"
	"net"
)

const SECRET = "7#u%AjH63p49"

func main() {
	listen, err := net.ListenTCP("tcp", &net.TCPAddr{net.ParseIP(""), 9090, ""})
	if err != nil {
		fmt.Println("监听端口失败:", err.Error())
		return
	}
	fmt.Println("已初始化连接，等待客户端连接...")
	Server(listen)
}

func Server(listen *net.TCPListener) {
	for {
		conn, err := listen.AcceptTCP()
		if err != nil {
			fmt.Println("接受客户端连接异常:", err.Error())
			continue
		}
		fmt.Println("客户端连接来自:", conn.RemoteAddr().String())
		defer conn.Close()
		go func() {
			data := make([]byte, 4096)
			for {
				i, err := conn.Read(data)
				fmt.Println("客户端发来数据:", string(data[0:i]))
				if err != nil {
					fmt.Println("读取客户端数据错误:", err.Error())
					break
				}
				conn.Write(data[:i])
			}

		}()
	}
}
