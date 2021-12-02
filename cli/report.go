package cli

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"

	"github.com/avarabyeu/goRP/gorp"
)

const logsBatchSize = 5

var (
	reportCommand = &cli.Command{
		Name:        "report",
		Usage:       "Reports input to report portal",
		Subcommands: cli.Commands{reportTest2JsonCommand},
	}

	reportTest2JsonCommand = &cli.Command{
		Name:  "test2json",
		Usage: "Input format: test2json",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:    "file",
				Aliases: []string{"f"},
				Usage:   "File Name",
				EnvVars: []string{"FILE"},
			},
			&cli.StringFlag{
				Name:    "launchName",
				Aliases: []string{"ln"},
				Usage:   "Launch Name",
				EnvVars: []string{"LAUNCH_NAME"},
				Value:   "gorp launch",
			},
		},
		Action: reportTest2json,
	}
)

func reportTest2json(c *cli.Context) error {
	rpClient, err := buildClient(c)
	if err != nil {
		return err
	}
	input := make(chan *testEvent)

	// run in separate goroutine
	launchNameArg := c.String("launchName")
	rep := newReporter(rpClient, launchNameArg, input)

	wg := &sync.WaitGroup{}
	go func() {
		wg.Add(1)
		defer wg.Done()
		rep.receive()
	}()
	defer wg.Wait()

	defer close(input)

	var reader io.Reader
	if fileName := c.String("file"); fileName != "" {
		f, fErr := os.Open(fileName)
		if fErr != nil {
			return fErr
		}
		defer func() {
			if cErr := f.Close(); cErr != nil {
				logrus.Error(cErr)
			}
		}()
		reader = f
	} else {
		reader = os.Stdin
	}

	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		data := scanner.Text()

		var ev testEvent
		if err := json.Unmarshal([]byte(data), &ev); err != nil {
			logrus.Error(err)
			return err
		}
		input <- &ev
	}
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
	input         <-chan *testEvent
	client        *gorp.Client
	launchName    string
	launchID      string
	launchOnce    sync.Once
	tests         map[string]string
	logs          []*gorp.SaveLogRQ
	logsBatchSize int
}

func newReporter(client *gorp.Client, launchName string, input <-chan *testEvent) *reporter {
	return &reporter{
		input:         input,
		launchName:    launchName,
		client:        client,
		launchOnce:    sync.Once{},
		tests:         map[string]string{},
		logs:          []*gorp.SaveLogRQ{},
		logsBatchSize: logsBatchSize,
	}
}

func (r *reporter) receive() {
	for ev := range r.input {
		if err := r.startLaunch(ev); err != nil {
			logrus.Error(err)
		}

		var err error
		switch ev.Action {
		case "run":
			_, err = r.startTest(ev)
		case "output":
			err = r.log(ev)
		case "pass":
			if ev.Test == "" {
				err = r.finishLaunch(ev, gorp.Statuses.Passed)
			} else {
				err = r.finishTest(ev, gorp.Statuses.Passed)
			}
		case "fail":
			if ev.Test == "" {
				err = r.finishLaunch(ev, gorp.Statuses.Failed)
			} else {
				err = r.finishTest(ev, gorp.Statuses.Failed)
			}
		}
		if err != nil {
			logrus.Fatal(err)
		}
	}
}

func (r *reporter) startTest(ev *testEvent) (string, error) {
	fmt.Println(ev.Time)
	testID := r.getTestName(ev)
	rs, err := r.client.StartTest(&gorp.StartTestRQ{
		StartRQ: gorp.StartRQ{
			Name:      ev.Test,
			StartTime: gorp.NewTimestamp(time.Now()),
		},
		LaunchID:   r.launchID,
		HasStats:   "true",
		UniqueID:   testID,
		CodeRef:    testID,
		TestCaseID: testID,
		Type:       gorp.TestItemTypes.Test,
		Retry:      false,
	})
	if err != nil {
		return "", err
	}
	r.tests[testID] = rs.ID
	return rs.ID, nil
}

func (r *reporter) log(ev *testEvent) error {
	if ev.Output == "" {
		return nil
	}
	testName := r.getTestName(ev)
	testID := r.tests[testName]

	// if output starts from tab
	if strings.HasPrefix(strings.TrimLeft(ev.Output, " "), "\t") && len(r.logs) > 0 {
		lastLog := r.logs[len(r.logs)-1]
		lastLog.Message = lastLog.Message + "\n" + ev.Output
		lastLog.Level = gorp.LogLevelError
		return nil
	}

	rq := &gorp.SaveLogRQ{
		ItemID:     testID,
		LaunchUUID: r.launchID,
		Level:      gorp.LogLevelInfo,
		LogTime:    gorp.NewTimestamp(time.Now()),
		Message:    ev.Output,
	}
	r.logs = append(r.logs, rq)
	if len(r.logs) >= r.logsBatchSize {
		_, err := r.client.SaveLogs(r.logs...)
		r.logs = []*gorp.SaveLogRQ{}
		if err != nil {
			return err
		}
	}
	return nil
}

func (r *reporter) getTestName(ev *testEvent) string {
	return fmt.Sprintf("%s/%s", ev.Package, ev.Test)
}

func (r *reporter) startLaunch(ev *testEvent) error {
	var err error
	r.launchOnce.Do(func() {
		var launch *gorp.EntryCreatedRS
		launch, err = r.client.StartLaunch(&gorp.StartLaunchRQ{
			StartRQ: gorp.StartRQ{
				Name:      r.launchName,
				StartTime: gorp.NewTimestamp(time.Now()),
			},
			Mode: gorp.LaunchModes.Default,
		})
		if err != nil {
			return
		}
		r.launchID = launch.ID
	})
	return err
}

func (r *reporter) finishLaunch(ev *testEvent, status gorp.Status) error {
	_, err := r.client.FinishLaunch(r.launchID, &gorp.FinishExecutionRQ{
		Status:  status,
		EndTime: gorp.NewTimestamp(time.Now()),
	})
	return err
}

func (r *reporter) finishTest(ev *testEvent, status gorp.Status) error {
	testName := r.getTestName(ev)
	testID := r.tests[testName]

	_, err := r.client.FinishTest(testID, &gorp.FinishTestRQ{
		FinishExecutionRQ: gorp.FinishExecutionRQ{
			EndTime: gorp.NewTimestamp(time.Now()),
			Status:  status,
		},
		LaunchUUID: r.launchID,
	})
	return err
}
