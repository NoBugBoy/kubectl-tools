package tcp

import (
	"fmt"
	"net"
)

type Socket struct {
	Conn net.Conn
}

func OpenSocket() (*Socket,chan int) {
	created := make(chan int,1)
	s := &Socket{
		Conn: nil,
	}
	netListen, err := net.Listen("tcp", ":19675")
	if err != nil{
		fmt.Println("close connection")
	}

	go func() {
		for {
			conn, err := netListen.Accept()
			if err != nil {
				continue
			}
			s.Conn = conn
			created <- 1
		}
	}()
	return s,created

}
