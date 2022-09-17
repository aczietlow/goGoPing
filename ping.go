package main

import (
	"fmt"
	"github.com/aczietlow/gogoping/cli"
	"github.com/aczietlow/gogoping/net"
	"math"
	"time"
)

func main() {
	terminal := cli.NewTerminal()
	options := terminal.Args.Options
	// Ensure that we gracefully restore the terminal connection.
	defer terminal.Restore()

	netClient := net.NewClient()
	// Ensure that as we exit gracefully we close the connection.
	defer netClient.Close()

	url := terminal.Args.Arg
	ip4 := netClient.ResolveIpAddress(url)

	fmt.Printf("PING %v (%v):  56 data bytes\r\n", url, ip4)

	// Grrr nothing likes to accept float32!!! Should probably fix this.
	// Gives us rounded to 2 decimal places and converts to milliseconds.
	wait := math.Round(float64(options.Wait*100)) * 10
	ms := time.Duration(wait)
	for i := 0; i < options.Count; i++ {
		time.Sleep(ms * time.Millisecond)
		netClient.Ping(ip4)
	}
	fmt.Printf("flag is %T\r\n", int64(wait))
	fmt.Printf("Fin\r\n")
}
