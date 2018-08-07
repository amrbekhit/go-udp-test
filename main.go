package main

import (
	"flag"
	"fmt"
	"os"
)

func main() {
	mode := flag.String("mode", "client", "Either 'server' or 'client'")
	address := flag.String("address", "", "The address to connect/listen to")

	flag.Parse()

	switch *mode {
	case "server":
		runServer(*address)
	case "client":
		runClient(*address)
	default:
		fmt.Printf("Invalid mode")
		os.Exit(1)
	}
}
