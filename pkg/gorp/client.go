package gorp

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

// Client is ReportPortal REST API Client
type Client struct {
	project string
	http    *resty.Client
}

// NewClient creates new instance of Client
// host - server hostname
// project - name of the project
// apiKey - User Token (see user profile page)
func NewClient(host, project, apiKey string) *Client {
	http := resty.New().
		// SetDebug(true).
		SetBaseURL(host).
		SetAuthToken(apiKey).
		OnAfterResponse(func(client *resty.Client, rs *resty.Response) error {
			//nolint:gomnd // 4xx errors
			if (rs.StatusCode() / 100) >= 4 {
				return fmt.Errorf("status code error: %d\n%s", rs.StatusCode(), rs.String())
			}

			return nil
		})

	return &Client{
		project: project,
		http:    http,
	}
}

// StartLaunch starts new launch in RP
func (c *Client) StartLaunch(launch *StartLaunchRQ) (*EntryCreatedRS, error) {
	return c.startLaunch(launch)
}

// StartLaunchRaw starts new launch in RP with body in form of bytes buffer
func (c *Client) StartLaunchRaw(body json.RawMessage) (*EntryCreatedRS, error) {
	return c.startLaunch(body)
}

// StartLaunch starts new launch in RP
func (c *Client) startLaunch(body interface{}) (*EntryCreatedRS, error) {
	var rs EntryCreatedRS
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetBody(body).
		SetResult(&rs).
		Post("/api/v2/{project}/launch")

	return &rs, err
}

// FinishLaunch finishes launch in RP
func (c *Client) FinishLaunch(id string, launch *FinishExecutionRQ) (*FinishLaunchRS, error) {
	return c.finishLaunch(id, launch)
}

// FinishLaunchRaw finishes launch in RP with body in form of bytes buffer
func (c *Client) FinishLaunchRaw(id string, body json.RawMessage) (*FinishLaunchRS, error) {
	return c.finishLaunch(id, body)
}

// FinishLaunch finishes launch in RP
func (c *Client) finishLaunch(id string, body interface{}) (*FinishLaunchRS, error) {
	var rs FinishLaunchRS
	_, err := c.http.R().
		SetPathParams(map[string]string{
			"project":  c.project,
			"launchId": id,
		}).
		SetBody(body).
		SetResult(&rs).
		Put("/api/v2/{project}/launch/{launchId}/finish")

	return &rs, err
}

// StopLaunch forces finishing launch
func (c *Client) StopLaunch(id string) (*MsgRS, error) {
	var rs MsgRS
	_, err := c.http.R().
		SetPathParams(map[string]string{
			"project":  c.project,
			"launchId": id,
		}).
		SetBody(&FinishExecutionRQ{
			EndTime: NewTimestamp(time.Now()),
			Status:  Statuses.Stopped,
		}).
		SetResult(&rs).
		Put("/api/v2/{project}/launch/{launchId}/stop")

	return &rs, err
}

// StartTest starts new test in RP
func (c *Client) StartTest(item *StartTestRQ) (*EntryCreatedRS, error) {
	return c.startTest(item)
}

// StartTestRaw starts new test in RP accepting request body as array of bytes
func (c *Client) StartTestRaw(body json.RawMessage) (*EntryCreatedRS, error) {
	return c.startTest(body)
}

// startTest starts new test in RP
func (c *Client) startTest(body interface{}) (*EntryCreatedRS, error) {
	var rs EntryCreatedRS
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetBody(body).
		SetResult(&rs).
		Post("/api/v2/{project}/item/")

	return &rs, err
}

// startChildTest starts new test in RP
func (c *Client) startChildTest(parent string, body interface{}) (*EntryCreatedRS, error) {
	var rs EntryCreatedRS
	_, err := c.http.R().
		SetPathParams(map[string]string{
			"project": c.project,
			"itemId":  parent,
		}).
		SetBody(body).
		SetResult(&rs).
		Post("/api/v2/{project}/item/{itemId}")

	return &rs, err
}

// StartChildTest starts new test in RP
func (c *Client) StartChildTest(parent string, item *StartTestRQ) (*EntryCreatedRS, error) {
	return c.startChildTest(parent, item)
}

// StartChildTestRaw starts new test in RP accepting request body as array of bytes
func (c *Client) StartChildTestRaw(parent string, body json.RawMessage) (*EntryCreatedRS, error) {
	return c.startChildTest(parent, body)
}

// FinishTest finishes test in RP
func (c *Client) FinishTest(id string, rq *FinishTestRQ) (*MsgRS, error) {
	return c.finishTest(id, rq)
}

// FinishTestRaw finishes test in RP accepting body as array of bytes
func (c *Client) FinishTestRaw(id string, body json.RawMessage) (*MsgRS, error) {
	return c.finishTest(id, body)
}

// finishTest finishes test in RP
func (c *Client) finishTest(id string, body interface{}) (*MsgRS, error) {
	var rs MsgRS
	_, err := c.http.R().
		SetPathParams(map[string]string{
			"project": c.project,
			"itemId":  id,
		}).
		SetBody(body).
		SetResult(&rs).
		Put("/api/v2/{project}/item/{itemId}")
	return &rs, err
}

// SaveLog attaches log in RP
func (c *Client) SaveLog(log *SaveLogRQ) (*EntryCreatedRS, error) {
	var rs EntryCreatedRS
	_, err := c.http.R().
		SetPathParams(map[string]string{
			"project": c.project,
		}).
		SetBody(log).
		SetResult(&rs).
		Post("/api/v2/{project}/log")
	return &rs, err
}

// SaveLogs saves logs as batch request
func (c *Client) SaveLogs(logs ...*SaveLogRQ) (*EntryCreatedRS, error) {
	return c.SaveLogMultipart(logs, nil)
}

// SaveLogMultipart saves a batch of logs in RP, along with any associated files (if any).
//
// Example usage:
//
// f, _ := os.Open("someFile.txt")
//
//	logs := []*SaveLogRQ{{
//		    File: FileAttachment{
//		        // note that this value must present in 'files' map as key (see below)
//		        Name: "fileAttachment.txt",
//		    },
//		    LaunchUUID: launchID,
//		    ItemID:     itemID,
//		    Level:      gorp.LogLevelError,
//		    LogTime:    NewTimestamp(time.Now()),
//		    Message:    "Important message!",
//			}}
//
//	files := map[string]*os.File{
//				"fileAttachment.txt": f, // key must match the FileAttachment.Name field
//			}
//
// resp, err := client.SaveLogMultipart(log, files)
func (c *Client) SaveLogMultipart(log []*SaveLogRQ, files []Multipart) (*EntryCreatedRS, error) {
	var bodyBuf bytes.Buffer
	err := json.NewEncoder(&bodyBuf).Encode(log)
	if err != nil {
		return nil, fmt.Errorf("unable to encode log payload: %w", err)
	}

	rq := c.http.R().
		SetPathParams(map[string]string{
			"project": c.project,
		})

	// JSON PAYLOAD PART
	rq.SetMultipartField("json_request_part", "", "application/json", &bodyBuf)

	// BINARY PART
	for _, v := range files {
		fileName, contentType, reader, lErr := v.Load()
		if lErr != nil {
			return nil, fmt.Errorf("unable to read multipart: %w", lErr)
		}
		if fileName == "" {
			return nil, errors.New("no file name is provided")
		}

		rq.SetMultipartField("file", fileName, contentType, reader)
	}

	var rs EntryCreatedRS
	_, err = rq.
		SetResult(&rs).
		Post("/api/v2/{project}/log")
	if err != nil {
		return nil, fmt.Errorf("unable to send log: %w", err)
	}
	return &rs, err
}

// GetLaunches retrieves latest launches
func (c *Client) GetLaunches() (*LaunchPage, error) {
	var launches LaunchPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

// GetLaunchesByFilter retrieves launches by filter
func (c *Client) GetLaunchesByFilter(filter map[string]string) (*LaunchPage, error) {
	var launches LaunchPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		SetQueryParams(filter).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

// GetLaunchesByFilterString retrieves launches by filter as string
func (c *Client) GetLaunchesByFilterString(filter string) (*LaunchPage, error) {
	var launches LaunchPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		SetQueryString(filter).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

// GetLaunchesByFilterName retrieves launches by filter name
func (c *Client) GetLaunchesByFilterName(name string) (*LaunchPage, error) {
	filter, err := c.GetFiltersByName(name)
	if err != nil {
		return nil, err
	}

	if filter.Page.Size < 1 || len(filter.Content) == 0 {
		return nil, fmt.Errorf("no filter %s found", name)
	}

	var launches LaunchPage
	params := ConvertToFilterParams(filter.Content[0])
	_, err = c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		SetQueryParams(params).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

// GetFiltersByName retrieves filter by its name
func (c *Client) GetFiltersByName(name string) (*FilterPage, error) {
	var filter FilterPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project, "name": name}).
		SetQueryParam("filter.eq.name", name).
		SetResult(&filter).
		Get("/api/v1/{project}/filter")
	return &filter, err
}

// MergeLaunches merge two launches
func (c *Client) MergeLaunches(rq *MergeLaunchesRQ) (*LaunchResource, error) {
	var rs LaunchResource
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetBody(rq).
		SetResult(&rs).
		Post("/api/v1/{project}/launch/merge")
	return &rs, err
}
