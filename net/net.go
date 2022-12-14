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
	DestIP     *network.IPAddr
}

func NewClient(options cli.Options, url string) *client {
	// Creates a new ICMP socket connection.
	c, err := icmp.ListenPacket("udp4", "0.0.0.0")
	if err != nil {
		log.Fatal(err)
	}

	err = c.IPv4PacketConn().SetTTL(options.TTL)
	if err != nil {
		log.Fatal(err)
	}

	// Setting the control message, passes IP header info back when we call ReadFrom() in a *controlMessage
	err = c.IPv4PacketConn().SetControlMessage(ipv4.FlagTTL, true)
	if err != nil {
		log.Fatal(err)
	}

	destIP := resolveIpAddress(url)

	return &client{
		Connection: c,
		DestIP:     destIP,
	}
}

func resolveIpAddress(url string) *network.IPAddr {
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

func (c *client) Ping(options cli.Options, seq int) {
	// @TODO Do we really need to create a new packet on every ping?
	wm := pingDatagram(options.Size, seq)

	wb, err := wm.Marshal(nil)
	if err != nil {
		log.Fatal(err)
	}

	start := time.Now()

	if _, err := c.Connection.WriteTo(wb, &network.UDPAddr{IP: network.ParseIP(c.DestIP.String())}); err != nil {
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
	t := time.Now()
	responseTime := t.Sub(start).Truncate(time.Microsecond)
	if err != nil {
		log.Fatal(err)
	}
	if cm != nil {
		ttl = cm.TTL
	}

	// Bitwise maths FTW, the seq is stored in 2 bytes. Then convert to int.
	replySeq := int(uint16(receivedPayload[6])<<8 | uint16(receivedPayload[7]))

	receivedICMPMessage, err := icmp.ParseMessage(ipv4.ICMPTypeEchoReply.Protocol(), receivedPayload[:numOfBytes])
	if err != nil {
		log.Fatal(err)
	}

	switch receivedICMPMessage.Type {
	case ipv4.ICMPTypeEchoReply:
		fmt.Printf("%v bytes received from %v: icmp_seq=%v ttl=%v responseTime=%s\r\n", numOfBytes, c.DestIP, replySeq, ttl, responseTime)
		if numOfBytes != options.Size {
			fmt.Printf("wrong total length %v instead of %v\r\n", numOfBytes, options.Size)
		}
	case ipv4.ICMPTypeDestinationUnreachable:
		fmt.Printf("The Host %s is unreachable\r\n", c.DestIP)
	case ipv4.ICMPTypeTimeExceeded:
		fmt.Printf("Host %s is slow\r\n", c.DestIP)
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
