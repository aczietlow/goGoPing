# GoGoPing

A rewrite of the ping application from C to go. Create a tool that is useful in debugging connections to the Internet.

## But Why Though?

1) Because I can
2) I thought it would be fun
3) Learning

## Requirements

* Take a hostname or ipaddress as input
  * Do a DNS Lookup
* Opens a socket
* When Ctl+c is pressed to exit, present the user with a report of aggregated statistics
* Support command line arguments and flags
  * `ping 127.0.0.1 -f -l 1400 -Fails`
  * -l
  * -t
  * -f

### Rules

1) Limit resources to the following go lang spec, wikipedia, networking RFCs, effective go, the go std library
2) Attempt to make use of concurrency
3) Half attempt writing a real app, and not a single giant spaghetti code mess

## Resources

* (Miro Board)[https://miro.com/app/board/uXjVPd_Mth8=/]

## TIL

### Struct fields exported or unexported

* Go's visibility flag are denoted by lowercase and capitalize letters
  * They're as Exported and unexported
* https://pkg.go.dev/golang.org/x/net/icmp exists
  * probably don't want to entirely reinvent this wheel

```go
package main
type animal struct {
	Cute bool
	food bool
	legs int
}
```

`animal.Cute` would be accessible to other packages
`animal.legs` & `animal.food` would not be accessible to other packages.

### defer statements

Defer functions are executed in LIFO (last in first out) order.

I'm using them as a way to ensure that claimed resources are restored once we're done with them.