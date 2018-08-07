package main

import (
	"bufio"
	"fmt"
	"net"
	"os"
	"strconv"
	"strings"

	"github.com/amrbekhit/go-udp-test/server"
)

func runServer(bindAddress string) {
	s := server.Server{}
	clients := make([]net.Addr, 0)

	err := s.Start(bindAddress)
	if err != nil {
		fmt.Println("Failed to start server:", err)
		return
	}
	defer s.Stop()

	txChan := make(chan server.AddressedString)
	defer close(txChan)

	// Receive
	go func() {
		for {
			rx, ok := <-s.RX
			if !ok {
				// The RX channel has been closed indicating we are quitting
				fmt.Println("Quitting RX")
				return
			}
			fmt.Printf("RX %v: %v\n", rx.Address, *rx.Message)
			clients = addClient(rx.Address, clients)
		}
	}()

	// Transmit
	go func() {
		for {
			tx, ok := <-txChan
			if !ok {
				fmt.Println("Quitting TX")
				return
			}
			s.TX <- tx
			fmt.Printf("TX %v: %v\n", tx.Address, *tx.Message)
		}
	}()

	// User interface
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		bytes, _, _ := reader.ReadLine()
		input := string(bytes)
		tokens := strings.Fields(input)
		if len(tokens) >= 1 {
			command := strings.ToUpper(tokens[0])
			switch command {
			case "Q", "QUIT":
				return
			case "C", "CLIENTS":
				// Print the client list
				if len(clients) > 0 {
					for i, c := range clients {
						fmt.Printf("%v: %v\n", i, c)
					}
				} else {
					fmt.Println("No clients yet")
				}
			case "S", "SEND":
				if len(tokens) < 3 {
					fmt.Println("To send: s clientid msg")
					continue
				}
				clientid, err := strconv.Atoi(tokens[1])
				if err != nil || clientid >= len(clients) {
					fmt.Println("Invalid client ID")
					continue
				}
				msg := strings.Join(tokens[2:], " ")
				s.TX <- server.AddressedString{Address: clients[clientid], Message: &msg}
			}
		}
	}
}

func addClient(c net.Addr, list []net.Addr) []net.Addr {
	for _, v := range list {
		if v.String() == c.String() {
			// this client already exists
			return list
		}
	}
	return append(list, c)
}
