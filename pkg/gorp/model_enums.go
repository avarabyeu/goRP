package gorp

// ENUMS
type (
	// LaunchMode - DEFAULT/DEBUG
	LaunchMode string
	// launchModeValuesType contains enum values for launch mode
	launchModeValuesType struct {
		Default LaunchMode
		Debug   LaunchMode
	}

	// MergeType is type of merge: BASIC or DEEP
	MergeType           string
	mergeTypeValuesType struct {
		Basic MergeType
		Deep  MergeType
	}

	// Status represents test item status
	Status           string
	statusValuesType struct {
		Passed      Status
		Failed      Status
		Stopped     Status
		Skipped     Status
		Interrupted Status
		Canceled    Status
		Info        Status
		Warn        Status
	}

	// TestItemType represents ENUM of test item types
	TestItemType           string
	testItemTypeValuesType struct {
		Suite        TestItemType
		Story        TestItemType
		Test         TestItemType
		Scenario     TestItemType
		Step         TestItemType
		BeforeClass  TestItemType
		BeforeGroups TestItemType
		BeforeMethod TestItemType
		BeforeSuite  TestItemType
		BeforeTest   TestItemType
		AfterClass   TestItemType
		AfterGroups  TestItemType
		AfterMethod  TestItemType
		AfterSuite   TestItemType
		AfterTest    TestItemType
	}
)

// LaunchModes is enum values for easy access
var LaunchModes = launchModeValuesType{
	Default: "DEFAULT",
	Debug:   "DEBUG",
}

// MergeTypes is enum values for easy access
var MergeTypes = mergeTypeValuesType{
	Deep:  "DEEP",
	Basic: "BASIC",
}

// Statuses is enum values for easy access
var Statuses = statusValuesType{
	Passed:      "PASSED",
	Failed:      "FAILED",
	Stopped:     "STOPPED",
	Skipped:     "SKIPPED",
	Interrupted: "INTERRUPTED",
	Canceled:    "CANCELLED", //nolint:misspell // defined as described on server end
	Info:        "INFO",
	Warn:        "WARN",
}

// TestItemTypes is enum values for easy access
var TestItemTypes = testItemTypeValuesType{
	Suite:        "SUITE",
	Story:        "STORY",
	Test:         "TEST",
	Scenario:     "SCENARIO",
	Step:         "STEP",
	BeforeClass:  "BEFORE_CLASS",
	BeforeGroups: "BEFORE_GROUPS",
	BeforeMethod: "BEFORE_METHOD",
	BeforeSuite:  "BEFORE_SUITE",
	BeforeTest:   "BEFORE_TEST",
	AfterClass:   "AFTER_CLASS",
	AfterGroups:  "AFTER_GROUPS",
	AfterMethod:  "AFTER_METHOD",
	AfterSuite:   "AFTER_SUITE",
	AfterTest:    "AFTER_TEST",
}
