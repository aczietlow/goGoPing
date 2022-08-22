package main

import (
	"fmt"
	"github.com/aczietlow/gogoping/cli"
)

func main() {
	terminal := cli.NewTerminal()

	defer terminal.Restore()

	fmt.Println("Enter an web address")
	url, err := terminal.Terminal.ReadLine()
	if err != nil {
		panic(err)
	}

	fmt.Printf("ping %s\n", url)
	fmt.Printf("ping %s\n", url)
	fmt.Printf("ping %s\n", url)

}
