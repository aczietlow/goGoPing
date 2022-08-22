package cli

import (
	"golang.org/x/term"
	"os"
)

type cli struct {
	State    *term.State
	Terminal *term.Terminal
}

func NewTerminal() *cli {
	oldState, err := term.MakeRaw(int(os.Stdin.Fd()))
	terminal := term.NewTerminal(os.Stdin, ">")
	if err != nil {
		panic(err)
	}
	return &cli{
		State:    oldState,
		Terminal: terminal,
	}
}

func (c *cli) Restore() {
	err := term.Restore(int(os.Stdin.Fd()), c.State)
	if err != nil {
		panic(err)
	}
}
