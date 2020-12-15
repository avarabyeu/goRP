package gorp

import (
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

func ExampleClient() {
	client := NewClient("xxx", "xxx", "xxx")

	launchUUID, _ := uuid.NewV4()
	launch, err := client.StartLaunch(&StartLaunchRQ{
		Mode: LaunchModes.Default,
		StartRQ: StartRQ{
			Name:        "gorp-test",
			UUID:        &launchUUID,
			StartTime:   Timestamp{Time: time.Now()},
			Description: "Demo Launch",
		},
	})
	checkErr(err, "unable to start launch")

	testUUID, _ := uuid.NewV4()
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

	log.Println("LAUNCHXXX")
	log.Println(launchUUID.String())

	_, err = client.SaveLog(&SaveLogRQ{
		LaunchUUID: launchUUID.String(),
		ItemID:     testUUID.String(),
		Level:      LogLevelInfo,
		LogTime:    Timestamp{time.Now()},
		Message:    "Log without binary",
	})
	checkErr(err, "unable to save log")

	file, _ := os.Open("../go.mod")
	_, err = client.SaveLogMultipart(&SaveLogRQ{
		LaunchUUID: launchUUID.String(),
		ItemID:     testUUID.String(),
		Level:      LogLevelInfo,
		Message:    "Log with binary",
	}, map[string]*os.File{"go.mod": file})

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
