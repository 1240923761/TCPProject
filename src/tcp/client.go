package main

import (
	"fmt"
	"io"
	"net"
	"os"
)

func SendFile(filePath string, fileSize int64, conn net.Conn) {
	f, err := os.Open(filePath)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer f.Close()
	var count int64
	for {
		buf := make([]byte, 2048)
		n, err := f.Read(buf)
		if err != nil && io.EOF == err {
			fmt.Println("文件传输完成")
			conn.Write([]byte("finish"))
			return
		}
		conn.Write(buf[:n])
		count += int64(n)
		sendPercent := float64(count) / float64(fileSize) * 100
		value := fmt.Sprintf("%.2f", sendPercent)

		//print the progress of update
		fmt.Println("uploadFile: " + value + `%`)
	}
}

func main() {
	fmt.Println("input file URL")
	var str string
	fmt.Scan(&str)
	fileInfo, err := os.Stat(str)
	if err != nil {
		fmt.Println(err)
		return
	}

	//create connection
	conn, err := net.Dial("tcp", ":8000")
	if err != nil {
		fmt.Println(err)
		return
	}
	defer conn.Close()
	fileName := fileInfo.Name()
	fileSize := fileInfo.Size()
	//send fileName to server
	conn.Write([]byte(fileName))
	buf := make([]byte, 2048)

	//read message from server
	n, err := conn.Read(buf)
	if err != nil {
		fmt.Println(err)
		return
	}
	recvData := string(buf[:n])
	if recvData == "ok" {
		//send file content
		SendFile(str, fileSize, conn)
	}
}
