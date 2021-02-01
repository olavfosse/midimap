package log

import (
	"fmt"

	"github.com/fossegrim/midimap/cm/helper"
	"github.com/fossegrim/midimap/driver"
	"github.com/fossegrim/midimap/lang/matcher"
	"github.com/fossegrim/midimap/usage"
	"gitlab.com/gomidi/midi"
	"gitlab.com/gomidi/midi/reader"
)

// Log corresponds to the log command modifier. args corresponds to
// the list of arguments listed on the command line after the log
// command modifier.
//
// For documentation about the log command modifier itself, consult
// the (to be written) manual.
// TODO: Write midimap(1) and/or midimap-log(1)
func Log(args []string) error {
	m, receivedMatcher, portNumber, err := parseArgs(args)
	if err != nil {
		return err
	}

	drv, err := driver.New()
	if err != nil {
		return err
	}
	defer drv.Close()

	ins, err := drv.Ins()
	if err != nil {
		return err
	}

	in, err := helper.GetInByPortNumber(ins, portNumber)
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
			if !receivedMatcher || helper.MatcherMatchesMessage(m, msg) {
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

func parseArgs(args []string) (m matcher.Matcher, receivedMatcher bool, portNumber uint64, err error) {
	if len(args) == 2 {
		m, err = matcher.Parse(args[1])
		if err != nil {
			return
		}
		receivedMatcher = true
	} else if len(args) != 1 {
		err = usage.Usage()
		return
	}
	portNumber, err = helper.ParsePortNumber(args[0])
	return
}
