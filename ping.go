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
	size := options.Size

	// TODO fix this. This is the size of the ICMP packet and IPheader (see Miro)
	// Include the 8 bytes from the header in when describing the total ICMP packet size.
	fmt.Printf("PING %v (%v):  %v data bytes\r\n", url, ip4, size)

	// Grrr nothing likes to accept float32!!! Should probably fix this.
	// Gives us rounded to 2 decimal places and converts to milliseconds.
	wait := math.Round(float64(options.Wait*100)) * 10
	ms := time.Duration(wait)
	for i := 0; i < options.Count; i++ {
		time.Sleep(ms * time.Millisecond)
		netClient.Ping(ip4, options, i+1)
	}
	// Debug'n shit
	//fmt.Printf("flag is %v\r\n", url)
	//fmt.Printf("Fin\r\n")
}
