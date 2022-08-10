package main

import (
	"log"

	tcp "github.com/fahrizalfarid/simple-tcp-udp-server/server/tcp"
	udp "github.com/fahrizalfarid/simple-tcp-udp-server/server/udp"
)

func main() {
	errsChan := make(chan error)

	go func() {
		tcpConn, err := tcp.New().Server()
		tcp.New().Handler(tcpConn)

		errsChan <- err
	}()

	go func() {
		udpConn, err := udp.New().Server()
		udp.New().Handler(udpConn)

		errsChan <- err
	}()

	for err := range errsChan {
		if err != nil {
			log.Fatal(err)
		}
	}
}
