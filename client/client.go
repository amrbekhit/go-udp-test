package client

import (
	"net"
	"sync"
)

// Client represents a UDP client
type Client struct {
	conn        *net.UDPConn
	RX, TX      chan *string
	finished    *sync.WaitGroup
	requestStop chan int
	udpAddr     *net.UDPAddr
}

// Start starts the UDP client
func (c *Client) Start(address string) error {
	var err error

	c.udpAddr, err = net.ResolveUDPAddr("udp", address)
	if err != nil {
		return err
	}

	c.conn, err = net.ListenUDP("udp", nil)
	if err != nil {
		return err
	}

	c.finished = new(sync.WaitGroup)
	c.finished.Add(2)
	c.requestStop = make(chan int)

	c.RX = make(chan *string)
	c.TX = make(chan *string)

	// Receive
	go func() {
	loop:
		for {
			buf := make([]byte, 65535)
			n, err := c.conn.Read(buf)
			if err != nil {
				select {
				case <-c.requestStop:
					// We've been requested to stop
					break loop
				default:
					// The server has gone away. Wait until
					// it comes back again
					continue
				}
			}
			msg := string(buf[:n])
			go func(msg *string) {
				c.RX <- msg
			}(&msg)
		}
		c.finished.Done()
	}()

	// Transmit
	go func() {
	loop:
		for {
			select {
			case msg := <-c.TX:
				c.conn.WriteTo([]byte(*msg), c.udpAddr)
			case <-c.requestStop:
				break loop
			}
		}
		c.finished.Done()
	}()

	return nil
}

// Stop stops the client and waits until all goroutines have quit
func (c *Client) Stop() {
	// Close the socket - this cancels the Receive goroutine Read method
	c.conn.Close()
	// Request both goroutines to quit
	c.requestStop <- 1
	c.requestStop <- 1
	// Wait for them to finish
	c.finished.Wait()
	close(c.RX)
	close(c.TX)
}
