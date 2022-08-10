package main

import (
	"fmt"
	"log"
	"net"
	"runtime"
	"sync"
	"time"
)

type client struct {
	totalMsg int
	delay    time.Duration
}

type clientInterface interface {
	runUdpClient(wg *sync.WaitGroup)
	runTcpClient(wg *sync.WaitGroup)
}

func (c *client) runUdpClient(wg *sync.WaitGroup) {
	defer wg.Done()

	var buf [1024]byte
	udpAddr := ":1800"

	str, err := net.ResolveUDPAddr("udp4", udpAddr)

	if err != nil {
		log.Fatal(err)
	}

	dial, err := net.DialUDP("udp4", nil, str)
	if err != nil {
		log.Fatal(err)
	}

	defer dial.Close()

	// read response from server if any
	go func() {
		for {
			n, err := dial.Read(buf[:])

			if err != nil {
				continue
			}

			fmt.Println(string(buf[:n]))
		}
	}()

	for i := 0; i < c.totalMsg; i++ {
		msg := fmt.Sprintf("Hy from udp client %v", i)
		dial.Write([]byte(msg))
		time.Sleep(c.delay * time.Millisecond)
	}

	_ = dial.Close()
}

func (c *client) runTcpClient(wg *sync.WaitGroup) {
	defer wg.Done()

	var buf [1024]byte
	tcpAddr := ":1801"

	str, err := net.ResolveTCPAddr("tcp", tcpAddr)
	if err != nil {
		log.Fatal(err)
	}

	dial, err := net.DialTCP("tcp", nil, str)
	if err != nil {
		log.Fatal(err)
	}

	defer dial.Close()

	// read response from server if any
	go func() {
		for {
			n, err := dial.Read(buf[:])

			if err != nil {
				continue
			}

			fmt.Println(string(buf[:n]))
		}
	}()

	for i := 0; i < c.totalMsg; i++ {
		msg := fmt.Sprintf("Hy from tcp client %v", i)
		dial.Write([]byte(msg))
		time.Sleep(c.delay * time.Millisecond)
	}

	_ = dial.Close()
}

func NewClient(totalMsg int, delay time.Duration) clientInterface {
	return &client{
		totalMsg: totalMsg,
		delay:    delay,
	}
}

func main() {
	runtime.GOMAXPROCS(12)

	var wg sync.WaitGroup
	now := time.Now()

	client := NewClient(5, 1000)

	// simulate there is 1000 client's connected at same time
	for i := 0; i < 1000; i++ {

		// if you're looking for speed use udp instead
		// tcp will eat a lot of memory and cpu
		wg.Add(1)
		go client.runTcpClient(&wg)

		wg.Add(1)
		go client.runUdpClient(&wg)
	}
	wg.Wait()

	fmt.Println("Done !!!", time.Since(now))
}
