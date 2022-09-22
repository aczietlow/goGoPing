package net

import (
	"fmt"
	"github.com/aczietlow/gogoping/cli"
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

func (c *client) Ping(targetIP *network.IPAddr, options cli.Options) {
	wm := pingDatagram(options.Size)

	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Connection.IPv4PacketConn().SetTTL(options.TTL)
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

	icmpMessage := make([]byte, 1500)

	// Shit I've tried to get the ip packet when reading the response and failed.
	// Stepping up to the packet connection, doesn't seem to give me what I need
	// I did learn that reading is clears the buffer, until a new request is received in the socket.

	//ipr := make([]byte, 1500)
	//c.Connection.IPv4PacketConn().ReadFrom(icmpMessage)
	//_, _, _ = c.Connection.ReadFrom(ipr)
	//if (runtime.GOOS == "darwin" || runtime.GOOS == "ios") && c.Connection.IPv4PacketConn() != nil {
	//numOfBytes, cm, _, _ := c.Connection.IPv4PacketConn().ReadFrom(icmpMessage)
	// numOfBytes, _, _ := c.Connection.IPv4PacketConn().PacketConn.ReadFrom(icmpMessage)
	//}

	// The OG statement
	numOfBytes, _, err := c.Connection.ReadFrom(icmpMessage)
	if err != nil {
		log.Fatal(err)
	}
	receivedMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), icmpMessage[:numOfBytes])
	if err != nil {
		log.Fatal(err)
	}

	switch receivedMessage.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("%v bytes received from %v: icmp_seq=0 ttl=101 time=34.905 ms\r\n", numOfBytes, targetIP)
		// @TODO detect if this is different.
		if numOfBytes != options.Size {
			fmt.Printf("wrong total length %v instead of %v\r\n", numOfBytes, options.Size)
		}
	case ipv4.ICMPTypeDestinationUnreachable:
		fmt.Printf("The Host %s is unreachable\r\n", targetIP)
	case ipv4.ICMPTypeTimeExceeded:
		fmt.Printf("Host %s is slow\r\n", targetIP)
	default:
		log.Printf("got %+v; want echo reply\r\n", receivedMessage)
	}
}

func pingDatagram(size int) icmp.Message {
	dataBytes := make([]byte, 0, size)
	dataBytes = append(dataBytes, 8)

	// 8 Bytes are reserved for the ICMP header data
	for len(dataBytes) < size-8 {
		l := len(dataBytes)
		// Because this is using hexadecimal, every 256 bytes added to the slice, start back at 0x00
		if l%256 == 0 {
			dataBytes = append(dataBytes, 00)
		} else {
			dataBytes = append(dataBytes, dataBytes[l-1]+01)
		}
	}

	return icmp.Message{
		Type: ipv4.ICMPTypeEcho, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1,
			Data: dataBytes,
		},
	}
}
