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

func (c *client) Ping(targetIP *network.IPAddr, options cli.Options, seq int) {
	wm := pingDatagram(options.Size, seq)

	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Connection.IPv4PacketConn().SetTTL(options.TTL)
	// Setting the control message, passes IP header info back when we call ReadFrom() in a *controlMessage
	err = c.Connection.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)
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

	receivedPayload := make([]byte, 1500)
	var cm *ipv4.ControlMessage
	var ttl, numOfBytes int
	numOfBytes, cm, _, err = c.Connection.IPv4PacketConn().ReadFrom(receivedPayload)
	if err != nil {
		log.Fatal(err)
	}

	receivedICMPMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), receivedPayload[:numOfBytes])
	if err != nil {
		log.Fatal(err)
	}

	if cm != nil {
		ttl = cm.TTL
	}
	//@TODO get the seq from the response packet
	switch receivedICMPMessage.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("%v bytes received from %v: icmp_seq=0 ttl=%v time=34.905 ms\r\n", numOfBytes, targetIP, ttl)
		if numOfBytes != options.Size {
			fmt.Printf("wrong total length %v instead of %v\r\n", numOfBytes, options.Size)
		}
	case ipv4.ICMPTypeDestinationUnreachable:
		fmt.Printf("The Host %s is unreachable\r\n", targetIP)
	case ipv4.ICMPTypeTimeExceeded:
		fmt.Printf("Host %s is slow\r\n", targetIP)
	default:
		log.Printf("got %+v; want echo reply\r\n", receivedICMPMessage)
	}
}

func pingDatagram(size int, seq int) icmp.Message {
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
			ID:   os.Getpid() & 0xffff,
			Seq:  seq,
			Data: dataBytes,
		},
	}
}
