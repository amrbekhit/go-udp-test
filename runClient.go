package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/amrbekhit/go-udp-test/client"
)

func runClient(address string) {
	c := client.Client{}

	err := c.Start(address)
	if err != nil {
		fmt.Println("Failed to start client:", err)
		return
	}
	defer c.Stop()

	txChan := make(chan *string)
	defer close(txChan)

	// Receive
	go func() {
		for {
			rx, ok := <-c.RX
			if !ok {
				fmt.Println("Quitting RX")
				return
			}
			fmt.Printf("RX: %v\n", *rx)
		}
	}()

	// Transmit + UI
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("> ")
		bytes, _, _ := reader.ReadLine()
		input := string(bytes)
		c.TX <- &input
	}
}
