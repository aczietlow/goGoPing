package net

import (
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"log"
	network "net"
	"os"
)

type net struct {
	Connection *icmp.PacketConn
}

func NewNet() *net {
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

	return &net{
		Connection: c,
	}
}

func ResolveIpAddress(url string) *network.IPAddr {
	ip, err := network.ResolveIPAddr("ip4", url)
	if err != nil {
		panic(err)
	}

	return ip
}

func (c *net) close() {

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
