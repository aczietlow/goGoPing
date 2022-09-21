package cli

import (
	"errors"
	"flag"
	"fmt"
	"strconv"
)

type Options struct {
	Count int
	Wait  float32
	Size  int
}

type arguments struct {
	Options Options
	Arg     string
}

type flagFloat32 float32

func initFlags() arguments {
	var flagFloat flagFloat32
	count := flag.Int("c", 2, "Stop after sending (and receiving) Count ECHO_RESPONSE packets.  If this option is not specified, ping will operate until interrupted.")
	// @TODO figure out if the super users part is enforced at the OS level or within Ping. & do that -
	// @TODO defining a custom flag var was a little extra for the actual needs here. Especially when we convert it to a floag64 later. Leaving this here for a hot minute as a lesson in interfaces and self-inflected pain.
	flag.Var(&flagFloat, "w", "Wait interval seconds between sending each packet. The default is to wait for one second between each packet normally, or not to wait in flood mode. Only super-user may set interval to values less 0.2 seconds.")

	// Size of data
	size := flag.Int("s", 56, "Specifies the number of data bytes to be sent. The default is 56, which translates into 64 ICMP data bytes when combined with the 8 bytes of ICMP header data.")

	// Called to parse all the flags.
	flag.Parse()

	if *size < 9 {
		*size = 9
	}

	// Arg() must be called after flag.Parse()
	arg := flag.Arg(0)

	// @TODO do weird BS time conversion on wait here to make it useful.
	return arguments{
		Options: Options{
			Count: *count,
			Wait:  float32(flagFloat),
			Size:  *size,
		},
		Arg: arg,
	}
}

//	String() string
//	Set(string) error

func (f *flagFloat32) String() string {
	return fmt.Sprintf("%v", *f)
}

func (f *flagFloat32) Set(s string) error {
	v, err := strconv.ParseFloat(s, 32)
	if err != nil {
		err = numError(err)
	}
	*f = flagFloat32(v)
	return err
}

func numError(err error) error {
	ne, ok := err.(*strconv.NumError)
	if !ok {
		return err
	}
	if ne.Err == strconv.ErrSyntax {
		return errors.New("parse error")
	}
	if ne.Err == strconv.ErrRange {
		return errors.New("value out of range")
	}
	return err
}
