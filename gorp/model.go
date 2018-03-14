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

//LaunchMode constants
const (
	LaunchModeDefault = "DEFAULT"
	LaunchModeDebug   = "DEBUG"
)
