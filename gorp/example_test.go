package gorp

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"
)

func ExampleClient() {
	client := NewClient("",
		"", "")

	launchUUID := uuid.New()
	launch, err := client.StartLaunch(&StartLaunchRQ{
		Mode: LaunchModes.Default,
		StartRQ: StartRQ{
			Name:        "gorp-test",
			UUID:        &launchUUID,
			StartTime:   NewTimestamp(time.Now()),
			Description: "Demo Launch",
		},
	})
	checkErr(err, "unable to start launch")

	testUUID := uuid.New()
	_, err = client.StartTest(&StartTestRQ{
		LaunchID: launch.ID,
		CodeRef:  "example_test.go",
		UniqueID: "another one unique ID",
		Retry:    false,
		Type:     TestItemTypes.Test,
		StartRQ: StartRQ{
			Name:      "Gorp Test",
			StartTime: Timestamp{time.Now()},
			UUID:      &testUUID,
		},
	})
	checkErr(err, "unable to start test")

	_, err = client.SaveLog(&SaveLogRQ{
		LaunchUUID: launchUUID.String(),
		ItemID:     testUUID.String(),
		Level:      LogLevelInfo,
		LogTime:    Timestamp{time.Now()},
		Message:    "Log without binary",
	})
	checkErr(err, "unable to save log")

	file1, _ := os.Open("../go.mod")
	file2, _ := os.Open("../go.sum")
	_, err = client.SaveLogMultipart([]*SaveLogRQ{
		{
			LaunchUUID: launchUUID.String(),
			ItemID:     testUUID.String(),
			Level:      LogLevelInfo,
			Message:    "Log with binary one",
			Attachment: FileAttachment{
				Name: "go.mod",
			},
		},
		{
			LaunchUUID: launchUUID.String(),
			ItemID:     testUUID.String(),
			Level:      LogLevelInfo,
			Message:    "Log with binary two",
			Attachment: FileAttachment{
				Name: "go.sum",
			},
		},
	}, map[string]*os.File{
		filepath.Base(file1.Name()): file1,
		filepath.Base(file2.Name()): file2,
	})

	checkErr(err, "unable to save log multipart")

	_, err = client.FinishTest(testUUID.String(), &FinishTestRQ{
		LaunchUUID: launchUUID.String(),
		FinishExecutionRQ: FinishExecutionRQ{
			EndTime: Timestamp{time.Now()},
			Status:  Statuses.Passed,
		},
	})
	checkErr(err, "unable to finish test")

	_, err = client.FinishLaunch(launchUUID.String(), &FinishExecutionRQ{
		Status:  Statuses.Passed,
		EndTime: Timestamp{time.Now()},
	})
	checkErr(err, "unable to finish launch")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
