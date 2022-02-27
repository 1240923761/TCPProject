package main

import (
	"fmt"
	"golang.org/x/time/rate"
	"net"
	"os"
	"runtime"
	"time"
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
	fmt.Println(addr + ": Transfer FileName from Client ：" + fileName)

	conn.Write([]byte("ok"))
	f, err := os.Create(fileName)
	if err != nil {
		fmt.Println(err)
		return
	}
	limiter := rate.NewLimiter(1024, 2048)
	for {
		if limiter.Allow() {
			buf := make([]byte, 2048)
			n, _ := conn.Read(buf)
			if string(buf) == "finish" {
				fmt.Println(addr + ": finish")
				runtime.Goexit()
			}
			fmt.Println(string(buf[:n]))
			f.Write(buf[:n])
			fmt.Printf("Ok  %s\n", time.Now().Format("2006-01-02 15:04:05.000"))
		} else {
			fmt.Printf("Err %s\n", time.Now().Format("2006-01-02 15:04:05.000"))
			time.Sleep(3 * time.Second)
		}

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
