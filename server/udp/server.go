package server

import (
	"fmt"
	"log"
	"net"
	"sync"
)

type server struct {
	addrUdp   string
	msgBuffer []byte
	mutex     sync.RWMutex
	network   string
}

type udpServerInterface interface {
	Server() (*net.UDPConn, error)
	Handler(*net.UDPConn)
}

func (s *server) Handler(udp *net.UDPConn) {
	for {
		var buff = s.msgBuffer
		n, addr, err := udp.ReadFromUDP(buff[:])

		if err != nil {
			continue
		}

		// do something here, such as insert data to db
		s.mutex.Lock()
		msg := fmt.Sprintf("udp message from %v, msg : %v", addr, string(buff[:n]))
		fmt.Println(msg)
		s.mutex.Unlock()

		// udp.Close()
	}
}

func (s *server) Server() (*net.UDPConn, error) {
	str, err := net.ResolveUDPAddr(s.network, s.addrUdp)

	if err != nil {
		log.Fatal(err)
	}

	udp, errListen := net.ListenUDP(s.network, str)

	if errListen != nil {
		log.Fatal(errListen)
	}

	fmt.Println("UDP Listening", udp.LocalAddr().String())
	return udp, errListen
}

func New() udpServerInterface {
	return &server{
		addrUdp:   ":1800",
		network:   "udp4",
		mutex:     sync.RWMutex{},
		msgBuffer: make([]byte, 1024),
	}
}
