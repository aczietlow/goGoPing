package main

import (
	"fmt"
	"github.com/aczietlow/gogoping/cli"
	"github.com/aczietlow/gogoping/net"
	"time"
)

func main() {
	terminal := cli.NewTerminal()
	// Ensure that we gracefully restore the terminal connection.
	defer terminal.Restore()

	netClient := net.NewClient()
	// Ensure that as we exit gracefully we close the connection.
	defer netClient.Close()

	const url = "google.com"
	ip4 := netClient.ResolveIpAddress(url)

	fmt.Printf("PING %v (%v):  56 data bytes\r\n", url, ip4)

	for i := 0; i < 2; i++ {
		time.Sleep(1 * time.Second)
		netClient.Ping(ip4)
	}

	fmt.Printf("Fin\r\n")
}
