package gorp

import (
	"strconv"
	"strings"
	"time"
)

type (
	//LaunchMode - DEFAULT/DEBUG
	LaunchMode string

	//Response is a representation of server response
	Response struct {
		//Page is a slice of data returned by server
		Page struct {
			Number        int `json:"number,omitempty"`
			Size          int `json:"size,omitempty"`
			TotalElements int `json:"totalElements,omitempty"`
			TotalPages    int `json:"totalPages,omitempty"`
		} `json:"page,omitempty"`
	}

	//LaunchResource - GET Launch response model
	LaunchResource struct {
		ID                  string      `json:"id"`
		Name                string      `json:"name,omitempty"`
		Number              int64       `json:"number"`
		Description         string      `json:"description,omitempty"`
		StartTime           Timestamp   `json:"start_time,omitempty"`
		EndTime             Timestamp   `json:"end_time,omitempty"`
		Status              string      `json:"status,omitempty"`
		Tags                []string    `json:"tags,omitempty"`
		Mode                LaunchMode  `json:"mode,omitempty"`
		ApproximateDuration float32     `json:"approximateDuration,omitempty"`
		HasRetries          bool        `json:"hasRetries,omitempty"`
		Statistics          *Statistics `json:"statistics,omitempty"`
	}

	//FilterResource - GET Filter response model
	FilterResource struct {
		ID              string                `json:"id"`
		Name            string                `json:"name"`
		Type            string                `json:"type"`
		Owner           string                `json:"owner"`
		Entities        []*FilterEntity       `json:"entities"`
		SelectionParams *FilterSelectionParam `json:"selection_parameters,omitempty"`
	}

	//FilterEntity - One piece of filter
	FilterEntity struct {
		Field     string `json:"filtering_field"`
		Condition string `json:"condition"`
		Value     string `json:"value"`
	}

	//FilterPage - GET Filter response model
	FilterPage struct {
		Content []*FilterResource
		Response
	}

	//FilterSelectionParam - Describes filter ordering
	FilterSelectionParam struct {
		PageNumber int            `json:"page_number"`
		Orders     []*FilterOrder `json:"orders,omitempty"`
	}

	//FilterOrder - Describes ordering
	FilterOrder struct {
		SortingColumn string `json:"sorting_column"`
		Asc           bool   `json:"is_asc"`
	}

	//LaunchPage - GET Launch response model
	LaunchPage struct {
		Content []*LaunchResource
		Response
	}

	//Statistics is a execution stat details
	Statistics struct {
		Executions *struct {
			Total  string `json:"total,omitempty"`
			Passed string `json:"passed,omitempty"`
			Failed string `json:"failed,omitempty"`
		} `json:"executions,omitempty"`
		Defects *struct {
			Product       map[string]int `json:"product_bug,omitempty"`
			Automation    map[string]int `json:"automation_bug,omitempty"`
			System        map[string]int `json:"system_issue,omitempty"`
			ToInvestigate map[string]int `json:"to_investigate,omitempty"`
			NotDefect     map[string]int `json:"no_defect,omitempty"`
		} `json:"defects,omitempty"`
	}

	//Timestamp is a wrapper around Time to support
	//Epoch milliseconds
	Timestamp struct {
		time.Time
	}

	//MergeType is type of merge: BASIC or DEEP
	MergeType string

	//MergeLaunchesRQ payload representation
	MergeLaunchesRQ struct {
		Description             string    `json:"description,omitempty"`
		StartTime               Timestamp `json:"start_time,omitempty"`
		EndTime                 Timestamp `json:"end_time,omitempty"`
		ExtendSuitesDescription bool      `json:"extendSuitesDescription,omitempty"`
		Launches                []string  `json:"launches"`
		MergeType               MergeType `json:"merge_type,omitempty"`
		Mode                    string    `json:"mode,omitempty"`
		Tags                    []string  `json:"tags,omitempty"`
		Name                    string    `json:"name,omitempty"`
	}

	//StartRQ payload representation
	StartRQ struct {
		Name        string    `json:"name,omitempty"`
		Description string    `json:"description,omitempty"`
		Tags        []string  `json:"tags,omitempty"`
		StartTime   Timestamp `json:"start_time,omitempty"`
	}

	//StartLaunchRQ payload representation
	StartLaunchRQ struct {
		StartRQ
	}

	//FinishTestRQ payload representation
	FinishTestRQ struct {
		FinishExecutionRQ
		Retry bool `json:"retry,omitempty"`
	}

	//SaveLogRQ payload representation. Without attaches.
	SaveLogRQ struct {
		ItemID  string    `json:"item_id,omitempty"`
		LogTime Timestamp `json:"time,omitempty"`
		Message string    `json:"message,omitempty"`
		Level   string    `json:"level,omitempty"`
	}

	//StartTestRQ payload representation
	StartTestRQ struct {
		StartRQ
		Parameters []string `json:"parameters,omitempty"`
		UniqueID   string   `json:"unique_id,omitempty"`
		LaunchID   string   `json:"launch_id,omitempty"`
		Type       string   `json:"type,omitempty"`
	}

	//FinishExecutionRQ payload representation
	FinishExecutionRQ struct {
		EndTime     Timestamp `json:"end_time,omitempty"`
		Status      string    `json:"status,omitempty"`
		Description string    `json:"description,omitempty"`
		Tags        []string  `json:"tags,omitempty"`
	}

	//EntryCreatedRS payload
	EntryCreatedRS struct {
		ID string `json:"id,omitempty"`
	}

	//StartLaunchRS payload
	StartLaunchRS struct {
		ID string `json:"id,omitempty"`
	}

	//MsgRS successful operation response payload
	MsgRS struct {
		Msg string `json:"msg,omitempty"`
	}
)

//UnmarshalJSON converts Epoch milliseconds (timestamp) to appropriate object
func (rt *Timestamp) UnmarshalJSON(b []byte) (err error) {
	msInt, err := strconv.ParseInt(strings.Trim(string(b), "\""), 10, 64)
	if err != nil {
		return err
	}

	rt.Time = time.Unix(0, msInt*int64(time.Millisecond))
	return nil
}

//MarshalJSON converts Epoch milliseconds (timestamp) to appropriate object
func (rt *Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(rt.Time.In(time.UTC).UnixNano()/int64(time.Millisecond), 10)), nil
}

//Client constants
const (
	LaunchModeDefault = "DEFAULT"
	LaunchModeDebug   = "DEBUG"
	MergeTypeBasic    = "BASIC"
	MergeTypeDeep     = "DEEP"

	StatusStopped = "STOPPED"
	StatusPassed  = "PASSED"
	StatusFailed  = "FAILED"
)
