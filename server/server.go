package server

import (
	"net"
	"sync"
)

// AddressedString represents a message paired with a network address
type AddressedString struct {
	Message *string
	Address net.Addr
}

// Server represents the UDP server
type Server struct {
	conn        net.PacketConn
	TX, RX      chan AddressedString
	finished    *sync.WaitGroup
	requestStop chan int
}

// Start starts the UDP server
func (s *Server) Start(bindAddress string) error {
	var err error
	s.conn, err = net.ListenPacket("udp", bindAddress)
	if err != nil {
		return err
	}

	s.finished = new(sync.WaitGroup)
	s.finished.Add(2)
	s.requestStop = make(chan int)

	s.TX = make(chan AddressedString)
	s.RX = make(chan AddressedString)

	// Receive
	go func() {
		for {
			buf := make([]byte, 65535)
			n, addr, err := s.conn.ReadFrom(buf)
			if err != nil {
				// The socket has been closed, indicating we want to quit
				break
			}
			msg := string(buf[:n])
			go func(msg *string, addr net.Addr) {
				s.RX <- AddressedString{Message: msg, Address: addr}
			}(&msg, addr)
		}
		s.finished.Done()
	}()

	// Transmit
	go func() {
	loop:
		for {
			select {
			case addrMsg := <-s.TX:
				s.conn.WriteTo([]byte(*addrMsg.Message), addrMsg.Address)
			case <-s.requestStop:
				break loop
			}
		}
		s.finished.Done()
	}()

	return nil
}

// Stop stops the server and waits until all goroutines have quit
func (s *Server) Stop() {
	// Close the socket - this terminates the RX goroutine
	s.conn.Close()
	// Request the TX goroutine to quit
	s.requestStop <- 1
	// Wait for both goroutines to finish
	s.finished.Wait()
	close(s.RX)
	close(s.TX)
}
