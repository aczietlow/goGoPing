package main

import (
	"fmt"
	"github.com/aczietlow/gogoping/cli"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	"net"
	"os"
	"time"
)

func main() {
	terminal := cli.NewTerminal()

	defer terminal.Restore()

	//fmt.Println("Enter an web address")
	//url, err := terminal.Terminal.ReadLine()
	//if err != nil {
	//	panic(err)
	//}
	//
	//fmt.Printf("ping %s\n", url)
	//fmt.Printf("ping %s\n", url)
	foo := networkingMagic()
	fmt.Printf("%v", foo)

}

func networkingMagic() *net.IPAddr {
	// create a connection?
	// build a ICMP packet. Use type 3 as we just want an echo relay

	const url = "google.com"
	ip, err := net.ResolveIPAddr("ip4", url)
	if err != nil {
		panic(err)
	}

	for i := 0; i < 10; i++ {
		time.Sleep(1 * time.Second)
		Ping(ip)
	}

	return ip
}

func Ping(targetIP *net.IPAddr) {

	//const targetIP = "74.125.138.138"
	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}

	// wrap handler in closure.
	defer func(c *icmp.PacketConn) {
		err := c.Close()
		if err != nil {
			panic(err)
		}
	}(c)

	wm := icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}

	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := c.WriteTo(wb, &net.UDPAddr{IP: net.ParseIP(targetIP.String())}); err != nil {
		log.Fatal(err)
	}

	// Set deadline
	err = c.SetReadDeadline(time.Now().Add(time.Second * 1))
	if err != nil {
		fmt.Printf("Error on SetReadDeadline %v", err)
		panic(err)
	}

	rb := make([]byte, 1500)
	n, peer, err := c.ReadFrom(rb)
	if err != nil {
		log.Fatal(err)
	}
	rm, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:n])
	if err != nil {
		log.Fatal(err)
	}

	switch rm.Type {
	case ipv4.ICMPTypeEchoReply:
		log.Printf("got reflection from %v\x0A", peer)
		switch rm.Code {
		case 0:
			fmt.Printf("Reply received from %s\x0A", targetIP)
		case 3:
			fmt.Printf("The Host %s is unreachable\n", targetIP)
		case 11:
			fmt.Printf("Host %s is slow\n", targetIP)
		default:
			fmt.Printf("The Host %s is unreachable\n", targetIP)
		}
	default:
		log.Printf("got %+v; want echo reply", rm)
	}
}