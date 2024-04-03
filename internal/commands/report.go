package commands

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/urfave/cli/v2"

	gorp2 "github.com/reportportal/goRP/v5/pkg/gorp"
)

const logsBatchSize = 10

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
			&cli.StringSliceFlag{
				Name:    "attr",
				Aliases: []string{"a"},
				Usage:   "Launch attribute with format 'key:value'. Omitting a ':' separator will tag the launch with the value.",
			},
		},
		Action: reportTest2json,
	}
)

//nolint:nonamedreturns // for readability
func reportTest2json(c *cli.Context) (err error) {
	rpClient, err := buildClient(c)
	if err != nil {
		return err
	}
	input := make(chan *testEvent)

	// run in separate goroutine
	launchNameArg := c.String("launchName")
	attrArgs := c.StringSlice("attr")
	rep := newReporter(rpClient, launchNameArg, input, attrArgs...)

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		err = rep.receive()
	}()
	// wait for report to complete
	defer wg.Wait()

	defer close(input)

	var reader io.Reader
	if fileName := c.String("file"); fileName != "" {
		f, fErr := os.Open(filepath.Clean(fileName))
		if fErr != nil {
			return fErr
		}
		defer func() {
			if cErr := f.Close(); cErr != nil {
				slog.Error(cErr.Error())
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
			slog.Default().Error(err.Error())
			return err
		}
		input <- &ev
	}
	return nil
}

type testEvent struct {
	Time    time.Time `json:"time"` // encodes as an RFC3339-format string
	Action  string    `json:"action"`
	Package string    `json:"package"`
	Test    string    `json:"test"`
	Elapsed float64   `json:"elapsed"` // seconds
	Output  string    `json:"output"`
}

type reporter struct {
	input            <-chan *testEvent
	client           *gorp2.Client
	launchName       string
	launchID         string
	launchOnce       sync.Once
	launchAttributes []*gorp2.Attribute
	tests            map[string]string
	suites           map[string]string
	logs             []*gorp2.SaveLogRQ
	logsBatchSize    int
	waitQueue        sync.WaitGroup
}

func newReporter(client *gorp2.Client, launchName string, input <-chan *testEvent, launchAttrArgs ...string) *reporter {
	launchAttributes := make([]*gorp2.Attribute, 0, len(launchAttrArgs))
	for _, attr := range launchAttrArgs {
		// Separate the key:value pair. If `:` is not present, the entire string is considered the value and an empty key is used
		var p gorp2.Parameter
		if key, value, ok := strings.Cut(attr, ":"); ok {
			p.Key = key
			p.Value = value
		} else {
			p.Value = attr
		}
		launchAttributes = append(launchAttributes, &gorp2.Attribute{
			Parameter: p,
			System:    false,
		})
	}

	return &reporter{
		input:            input,
		launchName:       launchName,
		launchAttributes: launchAttributes,
		client:           client,
		launchOnce:       sync.Once{},
		tests:            map[string]string{},
		suites:           map[string]string{},
		logs:             []*gorp2.SaveLogRQ{},
		logsBatchSize:    logsBatchSize,
	}
}

func (r *reporter) reportEvent(ev *testEvent) error {
	var err error
	switch ev.Action {
	case "run":
		_, err = r.startTest(ev)
	case "output":
		r.log(ev)
	case "pass":
		err = r.finish(ev, gorp2.Statuses.Passed)
	case "fail":
		err = r.finish(ev, gorp2.Statuses.Failed)
	}
	return err
}

func (r *reporter) receive() error {
	prevEventTime := time.Now()
	for ev := range r.input {
		var err error
		startTime := ev.Time

		// start launch once
		// when first event comes
		r.launchOnce.Do(func() {
			if err = r.startLaunch(startTime); err != nil {
				slog.Error(err.Error())
			}
		})
		if err != nil {
			return err
		}

		// report event to ReportPortal
		err = r.reportEvent(ev)
		if err != nil {
			return err
		}

		// remember last's event time
		// for RP's finishLaunch
		prevEventTime = ev.Time
	}

	// make sure we flush all logs that are left
	r.flushLogs(true)
	// wait for requests to get sent
	r.waitQueue.Wait()

	// finish launch of started
	if r.launchID != "" {
		if err := r.finishLaunch(gorp2.Statuses.Passed, prevEventTime); err != nil {
			return err
		}
	}
	return nil
}

func (r *reporter) startSuite(ev *testEvent) (string, error) {
	rs, err := r.client.StartTest(&gorp2.StartTestRQ{
		StartRQ: gorp2.StartRQ{
			Name:      ev.Package,
			StartTime: gorp2.NewTimestamp(ev.Time),
		},
		LaunchID: r.launchID,
		HasStats: false,
		Type:     gorp2.TestItemTypes.Suite,
		Retry:    false,
	})
	if err != nil {
		return "", err
	}
	r.suites[ev.Package] = rs.ID
	return rs.ID, nil
}

func (r *reporter) startTest(ev *testEvent) (string, error) {
	testID := r.getTestName(ev)
	parentID, found := r.suites[ev.Package]
	if !found {
		var err error
		parentID, err = r.startSuite(ev)
		if err != nil {
			return "", err
		}
	}
	rs, err := r.client.StartChildTest(parentID, &gorp2.StartTestRQ{
		StartRQ: gorp2.StartRQ{
			Name:      ev.Test,
			StartTime: gorp2.NewTimestamp(ev.Time),
		},
		LaunchID:   r.launchID,
		HasStats:   true,
		UniqueID:   testID,
		CodeRef:    testID,
		TestCaseID: testID,
		Type:       gorp2.TestItemTypes.Test,
		Retry:      false,
	})
	if err != nil {
		return "", err
	}
	r.tests[testID] = rs.ID
	return rs.ID, nil
}

func (r *reporter) log(ev *testEvent) {
	if ev.Output == "" {
		return
	}
	testName := r.getTestName(ev)
	testID := r.tests[testName]

	// if output starts from tab
	if strings.HasPrefix(strings.TrimLeft(ev.Output, " "), "\t") && len(r.logs) > 0 {
		lastLog := r.logs[len(r.logs)-1]
		lastLog.Message = lastLog.Message + "\n" + ev.Output
		lastLog.Level = gorp2.LogLevelError
		return
	}

	rq := &gorp2.SaveLogRQ{
		ItemID:     testID,
		LaunchUUID: r.launchID,
		Level:      gorp2.LogLevelInfo,
		LogTime:    gorp2.NewTimestamp(ev.Time),
		Message:    ev.Output,
	}
	r.logs = append(r.logs, rq)
	r.flushLogs(false)
}

func (r *reporter) flushLogs(force bool) {
	if force || (len(r.logs) >= r.logsBatchSize) {
		batch := r.logs
		r.waitQueue.Add(1)
		go func(logs []*gorp2.SaveLogRQ) {
			defer r.waitQueue.Done()

			if _, err := r.client.SaveLogs(logs...); err != nil {
				slog.Error("unable to report logs", "error", err, "batch_length", len(logs))
			}
		}(batch)
		r.logs = []*gorp2.SaveLogRQ{}
	}
}

func (r *reporter) getTestName(ev *testEvent) string {
	return fmt.Sprintf("%s/%s", ev.Package, ev.Test)
}

func (r *reporter) startLaunch(startTime time.Time) error {
	var launch *gorp2.EntryCreatedRS
	launch, err := r.client.StartLaunch(&gorp2.StartLaunchRQ{
		StartRQ: gorp2.StartRQ{
			Name:       r.launchName,
			StartTime:  gorp2.NewTimestamp(startTime),
			Attributes: r.launchAttributes,
		},
		Mode: gorp2.LaunchModes.Default,
	})
	if err != nil {
		return err
	}
	r.launchID = launch.ID
	return err
}

func (r *reporter) finishLaunch(status gorp2.Status, endTime time.Time) error {
	_, err := r.client.FinishLaunch(r.launchID, &gorp2.FinishExecutionRQ{
		Status:  status,
		EndTime: gorp2.NewTimestamp(endTime),
	})
	return err
}

func (r *reporter) finishTest(ev *testEvent, status gorp2.Status) error {
	testName := r.getTestName(ev)
	testID := r.tests[testName]

	_, err := r.client.FinishTest(testID, &gorp2.FinishTestRQ{
		FinishExecutionRQ: gorp2.FinishExecutionRQ{
			EndTime: gorp2.NewTimestamp(ev.Time),
			Status:  status,
		},
		LaunchUUID: r.launchID,
	})
	return err
}

func (r *reporter) finish(ev *testEvent, status gorp2.Status) error {
	var err error
	if ev.Test == "" {
		err = r.finishSuite(ev, status)
	} else {
		err = r.finishTest(ev, status)
	}
	return err
}

func (r *reporter) finishSuite(ev *testEvent, status gorp2.Status) error {
	suiteID := r.suites[ev.Package]

	_, err := r.client.FinishTest(suiteID, &gorp2.FinishTestRQ{
		FinishExecutionRQ: gorp2.FinishExecutionRQ{
			EndTime: gorp2.NewTimestamp(ev.Time),
			Status:  status,
		},
		LaunchUUID: r.launchID,
	})
	return err
}
