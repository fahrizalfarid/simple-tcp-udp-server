package tcp

import (
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"time"
)

type server struct {
	addrTcp   string
	msgBuffer []byte
	mutex     sync.Mutex
	network   string
}

type tcpServerInterface interface {
	Server() (*net.TCPListener, error)
	Handler(*net.TCPListener)
}

func (s *server) handler(conn *net.TCPConn) {
	defer conn.Close()

	for {
		var buff = s.msgBuffer
		n, err := conn.Read(buff)

		if err == io.EOF {
			return
		}

		if err != nil {
			return
		}

		// do something here, such as insert data to db
		s.mutex.Lock()
		msg := fmt.Sprintf("tcp message from %v, msg : %v", conn.RemoteAddr(), string(buff[:n]))
		fmt.Println(msg)
		s.mutex.Unlock()
	}
}

func (s *server) Handler(listener *net.TCPListener) {
	for {
		conn, err := listener.AcceptTCP()

		if err != nil {
			log.Println(err)
			return
		}

		err = conn.SetKeepAlive(true)
		if err != nil {
			log.Fatal(err)
		}

		err = conn.SetKeepAlivePeriod(30 * time.Second)
		if err != nil {
			log.Fatal(err)
		}

		go s.handler(conn)

	}
}

func (s *server) Server() (*net.TCPListener, error) {
	str, err := net.ResolveTCPAddr(s.network, s.addrTcp)
	if err != nil {
		log.Fatal(err)
	}

	listen, err := net.ListenTCP(s.network, str)

	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("TCP Listening", listen.Addr().String())
	return listen, err
}

func New() tcpServerInterface {
	return &server{
		addrTcp:   ":1801",
		network:   "tcp",
		mutex:     sync.Mutex{},
		msgBuffer: make([]byte, 512),
	}
}
