package gorp

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gofrs/uuid"
)

func ExampleClient() {
	client := NewClient("https://reportportal.epam.com", "xxx", "xxx")

	launchUUID, _ := uuid.NewV4()
	launch, err := client.StartLaunch(&StartLaunchRQ{
		Mode: LaunchModeDefault,
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
		Type:     TypeItemTest,
		StartRQ: StartRQ{
			Name:      "Gorp Test",
			StartTime: Timestamp{time.Now()},
			UUID:      &testUUID,
		},
	})
	checkErr(err, "unable to start test")

	_, err = client.SaveLog(&SaveLogRQ{
		ItemID:  testUUID.String(),
		Level:   LogLevelInfo,
		LogTime: Timestamp{time.Now()},
		Message: "Log without binary",
	})
	checkErr(err, "unable to save log")

	file, _ := os.Open("../go.mod")
	_, err = client.SaveLogMultipart(&SaveLogRQ{
		ItemID:  testUUID.String(),
		Level:   LogLevelInfo,
		Message: "Log with binary",
	}, map[string]*os.File{
		"go.mod": file,
	})
	checkErr(err, "unable to save log multipart")

	_, err = client.FinishTest(testUUID.String(), &FinishTestRQ{
		FinishExecutionRQ: FinishExecutionRQ{
			EndTime: Timestamp{time.Now()},
			Status:  StatusPassed,
		},
	})
	checkErr(err, "unable to finish test")

	_, err = client.FinishLaunch(launchUUID.String(), &FinishExecutionRQ{
		Status:  StatusPassed,
		EndTime: Timestamp{time.Now()},
	})
	checkErr(err, "unable to finish launch")

	fmt.Println("OK")
}

func checkErr(err error, msg string) {
	if err != nil {
		log.Fatal(msg, err)
	}
}
