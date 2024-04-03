package gorp

import (
	"strconv"
	"strings"
	"time"

	"github.com/google/uuid"
)

type (
	// StartLaunchRQ payload representation
	StartLaunchRQ struct {
		StartRQ
		Mode    LaunchMode `json:"mode"`
		Rerun   bool       `json:"rerun,omitempty"`
		RerunOf *uuid.UUID `json:"rerunOf,omitempty"`
	}

	// FinishTestRQ payload representation
	FinishTestRQ struct {
		FinishExecutionRQ
		LaunchUUID string `json:"launchUuid,omitempty"`
		TestCaseID string `json:"testCaseId,omitempty"`
		Retry      bool   `json:"retry,omitempty"`
		RetryOf    string `json:"retryOf,omitempty"`
	}

	// SaveLogRQ payload representation.
	SaveLogRQ struct {
		LaunchUUID string     `json:"launchUuid,omitempty"`
		ItemID     string     `json:"itemUuid,omitempty"`
		LogTime    Timestamp  `json:"time,omitempty"`
		Message    string     `json:"message,omitempty"`
		Level      string     `json:"level,omitempty"`
		Attachment Attachment `json:"file,omitempty"`
	}

	// StartTestRQ payload representation
	StartTestRQ struct {
		StartRQ
		CodeRef    string       `json:"codeRef,omitempty"`
		Parameters []*Parameter `json:"parameters,omitempty"`
		UniqueID   string       `json:"uniqueId,omitempty"`
		TestCaseID string       `json:"testCaseId,omitempty"`
		LaunchID   string       `json:"launchUuid,omitempty"`
		Type       TestItemType `json:"type,omitempty"`
		Retry      bool         `json:"retry,omitempty"`
		HasStats   bool         `json:"hasStats,omitempty"`
	}

	// FinishExecutionRQ payload representation
	FinishExecutionRQ struct {
		EndTime     Timestamp    `json:"end_time,omitempty"`
		Status      Status       `json:"status,omitempty"`
		Description string       `json:"description,omitempty"`
		Attributes  []*Attribute `json:"attribute,omitempty"`
	}

	// EntryCreatedRS payload
	EntryCreatedRS struct {
		ID string `json:"id,omitempty"`
	}

	// StartLaunchRS payload
	StartLaunchRS struct {
		ID string `json:"id,omitempty"`
	}

	// MsgRS successful operation response payload
	MsgRS struct {
		Msg string `json:"msg,omitempty"`
	}

	// StartRQ payload representation
	StartRQ struct {
		UUID        *uuid.UUID   `json:"uuid,omitempty"`
		Name        string       `json:"name,omitempty"`
		Description string       `json:"description,omitempty"`
		Attributes  []*Attribute `json:"attributes,omitempty"`
		StartTime   Timestamp    `json:"start_time,omitempty"`
	}

	// Attribute represents ReportPortal's attribute
	Attribute struct {
		Parameter
		System bool `json:"system,omitempty"`
	}

	// Parameter represents key-value pair
	Parameter struct {
		Key   string `json:"key,omitempty"`
		Value string `json:"value,omitempty"`
	}

	// FinishLaunchRS is finish execution payload
	FinishLaunchRS struct {
		EntryCreatedRS
		Number int64 `json:"number,omitempty"`
	}
	// Timestamp is a wrapper around Time to support
	// Epoch milliseconds
	Timestamp struct {
		time.Time
	}

	// Attachment represents file attachment in log entries
	Attachment struct {
		Name string `json:"name,omitempty"`
	}
)

// UnmarshalJSON converts Epoch milliseconds (timestamp) to appropriate object
func (rt *Timestamp) UnmarshalJSON(b []byte) error {
	trimmed := strings.Trim(string(b), "\"")
	msInt, err := strconv.ParseInt(trimmed, 10, 64)
	if err != nil {
		dt, err := time.Parse(defaultDateTimeFormat, trimmed)
		if err != nil {
			return err
		}
		rt.Time = dt
		return nil
	}

	rt.Time = time.Unix(0, msInt*int64(time.Millisecond))
	return nil
}

// MarshalJSON converts Epoch milliseconds (timestamp) to appropriate object
func (rt *Timestamp) MarshalJSON() ([]byte, error) {
	return []byte(strconv.FormatInt(rt.Time.In(time.UTC).UnixNano()/int64(time.Millisecond), 10)), nil
}

// NewTimestamp creates Timestamp wrapper for time.Time
func NewTimestamp(t time.Time) Timestamp {
	return Timestamp{Time: t}
}
