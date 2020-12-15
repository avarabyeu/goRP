package gorp

const defaultDateTimeFormat = "2006-01-02T15:04:05.999-0700"

// Client constants
const (
	LogLevelDebug = "DEBUG"
	LogLevelInfo  = "INFO"
)

// PAYLOADS
type (
	// Response is a representation of server response
	Response struct {
		// Page is a slice of data returned by server
		Page struct {
			Number        int `json:"number,omitempty"`
			Size          int `json:"size,omitempty"`
			TotalElements int `json:"totalElements,omitempty"`
			TotalPages    int `json:"totalPages,omitempty"`
		} `json:"page,omitempty"`
	}

	// LaunchResource - GET Launch response model
	LaunchResource struct {
		ID                  string      `json:"id"`
		Name                string      `json:"name,omitempty"`
		Number              int64       `json:"number"`
		Description         string      `json:"description,omitempty"`
		StartTime           Timestamp   `json:"start_time,omitempty"`
		EndTime             Timestamp   `json:"end_time,omitempty"`
		Status              Status      `json:"status,omitempty"`
		Tags                []string    `json:"tags,omitempty"`
		Mode                LaunchMode  `json:"mode,omitempty"`
		ApproximateDuration float32     `json:"approximateDuration,omitempty"`
		HasRetries          bool        `json:"hasRetries,omitempty"`
		Statistics          *Statistics `json:"statistics,omitempty"`
	}

	// FilterResource - GET Filter response model
	FilterResource struct {
		ID              string                `json:"id"`
		Name            string                `json:"name"`
		Type            TestItemType          `json:"type"`
		Owner           string                `json:"owner"`
		Entities        []*FilterEntity       `json:"entities"`
		SelectionParams *FilterSelectionParam `json:"selection_parameters,omitempty"`
	}

	// FilterEntity - One piece of filter
	FilterEntity struct {
		Field     string `json:"filtering_field"`
		Condition string `json:"condition"`
		Value     string `json:"value"`
	}

	// FilterPage - GET Filter response model
	FilterPage struct {
		Content []*FilterResource
		Response
	}

	// FilterSelectionParam - Describes filter ordering
	FilterSelectionParam struct {
		PageNumber int            `json:"page_number"`
		Orders     []*FilterOrder `json:"orders,omitempty"`
	}

	// FilterOrder - Describes ordering
	FilterOrder struct {
		SortingColumn string `json:"sorting_column"`
		Asc           bool   `json:"is_asc"`
	}

	// LaunchPage - GET Launch response model
	LaunchPage struct {
		Content []*LaunchResource
		Response
	}

	// Statistics is a execution stat details
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

	// MergeLaunchesRQ payload representation
	MergeLaunchesRQ struct {
		Description             string     `json:"description,omitempty"`
		StartTime               *Timestamp `json:"start_time,omitempty"`
		EndTime                 *Timestamp `json:"end_time,omitempty"`
		ExtendSuitesDescription bool       `json:"extendSuitesDescription,omitempty"`
		Launches                []string   `json:"launches"`
		MergeType               MergeType  `json:"merge_type,omitempty"`
		Mode                    string     `json:"mode,omitempty"`
		Tags                    []string   `json:"tags,omitempty"`
		Name                    string     `json:"name,omitempty"`
	}
)
