package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
)

type server struct {
	addrTcp   string
	msgBuffer []byte
	mutex     sync.RWMutex
	network   string
}

type tcpServerInterface interface {
	Server() (*net.Listener, error)
	Handler(net.Listener)
}

func (s *server) Handler(listener net.Listener) {
	for {
		var buff = s.msgBuffer
		conn, err := listener.Accept()

		if err != nil {
			continue
		}

		go func() {
			for {
				n, err := conn.Read(buff)

				if err == io.EOF {
					continue
				}

				if err != nil {
					continue
				}

				// do something here, such as insert data to db
				s.mutex.Lock()
				msg := fmt.Sprintf("tcp message from %v, msg : %v", conn.RemoteAddr(), string(buff[:n]))
				fmt.Println(msg)
				s.mutex.Unlock()
			}
		}()
		// conn.Close()
	}
}

func (s *server) Server() (*net.Listener, error) {
	listen, err := net.Listen(s.network, s.addrTcp)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TCP Listening", listen.Addr().String())
	return &listen, err
}

func New() tcpServerInterface {
	return &server{
		addrTcp:   ":1801",
		network:   "tcp",
		mutex:     sync.RWMutex{},
		msgBuffer: make([]byte, 1024),
	}
}
