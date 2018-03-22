package main

import (
	"net"
	"time"
    "fmt"
)

func main() {
	service := ":8404"
	tcpAddr, err := net.ResolveTCPAddr("tcp4", service)
	checkError(err)
	listener, err := net.ListenTCP("tcp", tcpAddr)
	checkError(err)
    fmt.Println("begin listening") 
	for {
		conn, err := listener.Accept()
		if err != nil {
			continue
		}
		daytime := time.Now().Format("2006-01-02 15:04:05")
		conn.Write([]byte(daytime)) // don't care about return value
		conn.Close()                // we're finished with this client
	}
}
func checkError(err error) {
	if err != nil {
		panic(err)
	}
}
