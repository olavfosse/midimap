package main

import (
	"fmt"

	"github.com/fossegrim/midimap/lang/matcher"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
)

// logCommandModifier corresponds to the log command modifier. args corresponds
// to the list of arguments listed on the command line after the log command
// modifier.
//
// For documentation about the log command modifier itself, consult midimap(1).
func logCommandModifier(args []string) error {
	// Parse args
	var m matcher.Matcher
	var receivedMatcher bool
	switch len(args) {
	case 2:
		var err error
		m, err = matcher.Parse(args[1])
		if err != nil {
			return err
		}
	case 1:
	default:
		return errUsage
	}
	portNumber, err := parsePortNumber(args[0])
	if err != nil {
		return err
	}

	drv, err := newDriver()
	if err != nil {
		return err
	}
	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		return err
	}

	in, err := getInByPortNumber(ins, portNumber)
	if err != nil {
		return err
	}

	err = in.Open()
	if err != nil {
		return err
	}
	defer in.Close()

	rd := reader.New(
		reader.NoLogger(),
		reader.Each(func(pos *reader.Position, msg midi.Message) {
			if !receivedMatcher || matcherMatchesMessage(m, msg) {
				fmt.Println("---===---")
				fmt.Printf("%s\n", msg)
				fmt.Printf("status: %d\n", msg.Raw()[0])
				fmt.Printf("data1: %d\n", msg.Raw()[1])
				fmt.Printf("data2: %d\n", msg.Raw()[2])
			}
		}),
	)

	// I don't understand how and why this snippet works, but it does.
	// I asked for an explanation in https://github.com/vipul-sharma20/midi-macro/issues/1.
	// TODO: Read chapter 8 in the go programming language, which hopefully explains this.
	exit := make(chan string)

	go rd.ListenTo(in)

	for {
		select {
		case <-exit:
			return nil
		}
	}
}
