package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"github.com/avarabyeu/goRP/gorp"
	"io"
	"os"
	"strings"
	"time"

	"github.com/urfave/cli/v2"
)

var (
	reportCommand = &cli.Command{
		Name:        "report",
		Usage:       "Reports input to report portal",
		Subcommands: cli.Commands{reportTest2JsonCommand},
	}

	reportTest2JsonCommand = &cli.Command{
		Name:   "test2json",
		Usage:  "Input format: test2json",
		Action: reportTest2json,
	}
)

func reportTest2json(c *cli.Context) error {
	rpClient, err := buildClient(c)
	if err != nil {
		return err
	}
	reader := bufio.NewReader(os.Stdin)
	input := make(chan *testEvent)

	// run in separate goroutine
	go newReporter(rpClient, input).Report()

	defer close(input)

	for {
		var data string
		var readErr error
		data, readErr = reader.ReadString('\n')
		if len(data) == 0 && err == io.EOF {
			break
		}
		data = strings.TrimSuffix(data, "\n")
		var ev testEvent
		if err := json.Unmarshal([]byte(data), &ev); err != nil {
			return err
		}
		input <- &ev

		if readErr != nil {
			break
		}
	}
	if err != io.EOF {
		return err
	}
	fmt.Println(rpClient)
	return nil
}

type testEvent struct {
	Time    time.Time // encodes as an RFC3339-format string
	Action  string
	Package string
	Test    string
	Elapsed float64 // seconds
	Output  string
}

type reporter struct {
	input  <-chan *testEvent
	client *gorp.Client
}

func newReporter(client *gorp.Client, input <-chan *testEvent) *reporter {
	return &reporter{input: input, client: client}
}

func (r *reporter) Report() {
	for ev := range r.input {
		fmt.Println(ev)
	}
}
