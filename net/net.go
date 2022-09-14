package net

import (
	"fmt"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	network "net"
	"os"
	"time"
)

type client struct {
	Connection *icmp.PacketConn
}

func NewClient() *client {
	// Creates a new ICMP socket connection.
	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}

	// @TODO Need to figure out how this will work with the Net struct
	// Wrap handler in a closure.
	//defer func(c *icmp.PacketConn) {
	//	err := c.Close()
	//	if err != nil {
	//		panic(err)
	//	}
	//}(c)

	return &client{
		Connection: c,
	}
}

func (c client) ResolveIpAddress(url string) *network.IPAddr {
	ip, err := network.ResolveIPAddr("ip4", url)
	if err != nil {
		panic(err)
	}

	return ip
}

func (c *client) Close() {
	err := c.Connection.Close()
	if err != nil {
		panic(err)
	}
}

func (c *client) Ping(targetIP *network.IPAddr) {
	//const targetIP = "74.125.138.138"

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
	if _, err := c.Connection.WriteTo(wb, &network.UDPAddr{IP: network.ParseIP(targetIP.String())}); err != nil {
		log.Fatal(err)
	}

	// Set deadline
	err = c.Connection.SetReadDeadline(time.Now().Add(time.Second * 1))
	if err != nil {
		fmt.Printf("Error on SetReadDeadline %v", err)
		panic(err)
	}

	rb := make([]byte, 1500)
	numOfBytes, _, err := c.Connection.ReadFrom(rb)

	if err != nil {
		log.Fatal(err)
	}
	receivedMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), rb[:numOfBytes])
	if err != nil {
		log.Fatal(err)
	}

	// @TODO numOfBytes is the size of the ICMP message. Adding 20 bytes hard codes adding the ICMP packet and IP4 transport protocol.
	bytes := numOfBytes + 20
	switch receivedMessage.Type {
	case ipv4.ICMPTypeEchoReply:

		switch receivedMessage.Code {
		// Echo Ping Reply
		case 0:
			fmt.Printf("%v bytes received from %v: icmp_seq=0 ttl=56 time=34.905 ms\r\n", bytes, targetIP)
		case 3:
			fmt.Printf("The Host %s is unreachable\r\n", targetIP)
		case 11:
			fmt.Printf("Host %s is slow\r\n", targetIP)
		default:
			fmt.Printf("The Host %s is unreachable\r\n", targetIP)
		}
	default:
		log.Printf("got %+v; want echo reply\r\n", receivedMessage)
	}
}

func pingDatagram() icmp.Message {
	return icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: []byte("HELLO-R-U-THERE"),
		},
	}
}
