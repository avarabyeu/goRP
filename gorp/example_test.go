package gorp

import (
	"fmt"
	"github.com/gofrs/uuid"
	"log"
	"os"
	"time"
)

func ExampleClient() {
	client := NewClient("https://dev.reportportal.io", "default_personal", "xxxx")

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
	checkErr(err)

	testUUID, _ := uuid.NewV4()
	test, err := client.StartTest(&StartTestRQ{
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
	checkErr(err)
	fmt.Println(test)

	_, err = client.SaveLog(&SaveLogRQ{
		ItemID:  testUUID.String(),
		Level:   LogLevelInfo,
		LogTime: Timestamp{time.Now()},
		Message: "Log without binary",
	})
	checkErr(err)

	file, _ := os.Open("/Users/avarabyeu/work/sources/own/goRP/go.mod")
	_, err = client.SaveLogMultipart(&SaveLogRQ{
		ItemID:  testUUID.String(),
		Level:   LogLevelInfo,
		Message: "Log without binary",
	}, map[string]*os.File{
		"go.mod": file,
	})
	checkErr(err)

	finishTest, err := client.FinishTest(testUUID.String(), &FinishTestRQ{
		FinishExecutionRQ: FinishExecutionRQ{
			EndTime: Timestamp{time.Now()},
			Status:  StatusPassed,
		},
	})
	checkErr(err)

	fmt.Println(finishTest)

	finishLaunch, err := client.FinishLaunch(launchUUID.String(), &FinishExecutionRQ{
		Status:  StatusPassed,
		EndTime: Timestamp{time.Now()},
	})
	checkErr(err)

	fmt.Println(finishLaunch)
	fmt.Println("OK")

	// Output: OK
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
