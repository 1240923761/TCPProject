package main

import (
	"fmt"
	"net"
	"os"
	"runtime"
)

func Handler(conn net.Conn) {
	buf := make([]byte, 2048)
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	fileName := string(buf[:n])
	addr := conn.RemoteAddr().String()
	fmt.Println(addr + ": Transfer FileName from Client \n" + fileName)

	conn.Write([]byte("ok"))
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}

	for {
		buf := make([]byte, 2048)
		n, _ := conn.Read(buf)
		if string(buf[:n]) == "finish" {
			fmt.Println(addr + ": finish")
			runtime.Goexit()
		}
		f.Write(buf[:n])

	}
	defer conn.Close()
	defer f.Close()
}
func main() {
	listen, err := net.Listen("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer listen.Close()
	fmt.Println("server is running")
	for {
		conn, err := listen.Accept()
		if err != nil {
			fmt.Println(err)
			return
		}
		go Handler(conn)
	}
}
