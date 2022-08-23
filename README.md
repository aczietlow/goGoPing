# GoGoPing

A rewrite of the ping application from C to go

## But Why Though?

1) Because I can
2) I thought it would be fun
3) Learning

## Rules

1) Limit resources to the following go lang spec, wikipedia, networking RFCs, effective go, the go std library
2) Attempt to make use of concurrency
3) Half attempt writing a real app, and not a single giant spaghetti code mess

## Resources

* (Miro Board)[https://miro.com/app/board/uXjVPd_Mth8=/]

## TIL

* Go's visibility flag are denoted by lowercase and capitalize letters
* https://pkg.go.dev/golang.org/x/net/icmp exists
** probably don't want to entirely reinvent this wheel


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