package gorp

import (
	"fmt"
	"gopkg.in/resty.v1"
	"time"
)

//Client is ReportPortal REST API Client
type Client struct {
	project string
	http    *resty.Client
}

//NewClient creates new instance of Client
//host - server hostname
//project - name of the project
//uuid - User Token (see user profile page)
func NewClient(host, project, uuid string) *Client {
	http := resty.New().
		//SetDebug(true).
		SetHostURL(host).
		SetAuthToken(uuid).
		OnAfterResponse(func(client *resty.Client, rs *resty.Response) error {
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

//StartLaunch starts new launch in RP
func (c *Client) StartLaunch(launch *StartLaunchRQ) (*EntryCreatedRS, error) {
	var rs EntryCreatedRS
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetBody(launch).
		SetResult(&rs).
		Post("/api/v1/{project}/launch")
	return &rs, err
}

//FinishLaunch finishes launch in RP
func (c *Client) FinishLaunch(id string, launch *FinishExecutionRQ) (*MsgRS, error) {
	var rs MsgRS
	_, err := c.http.R().
		SetPathParams(map[string]string{
			"project":  c.project,
			"launchId": id,
		}).
		SetBody(launch).
		SetResult(&rs).
		Put("/api/v1/{project}/launch/{launchId}/finish")
	return &rs, err
}

//StopLaunch forces finishing launch
func (c *Client) StopLaunch(id string) (*MsgRS, error) {
	var rs MsgRS
	_, err := c.http.R().
		SetPathParams(map[string]string{
			"project":  c.project,
			"launchId": id,
		}).
		SetBody(&FinishExecutionRQ{
			EndTime: Timestamp{Time: time.Now()},
			Status:  StatusStopped,
		}).
		SetResult(&rs).
		Put("/api/v1/{project}/launch/{launchId}/stop")
	return &rs, err
}

//StartTest starts new test in RP
func (c *Client) StartTest(item *StartTestRQ) (*EntryCreatedRS, error) {
	var rs EntryCreatedRS
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetBody(item).
		SetResult(&rs).
		Post("/api/v1/{project}/item/")
	return &rs, err
}

//StartChildTest starts new test in RP
func (c *Client) StartChildTest(parent string, item *StartTestRQ) (*EntryCreatedRS, error) {
	var rs EntryCreatedRS
	_, err := c.http.R().
		SetPathParams(map[string]string{
			"project": c.project,
			"itemId":  parent,
		}).
		SetBody(item).
		SetResult(&rs).
		Post("/api/v1/{project}/item/{itemId}")
	return &rs, err
}

//FinishTest finishes test in RP
func (c *Client) FinishTest(id string, launch *FinishTestRQ) (*MsgRS, error) {
	var rs MsgRS
	_, err := c.http.R().
		SetPathParams(map[string]string{
			"project": c.project,
			"itemId":  id,
		}).
		SetBody(launch).
		SetResult(&rs).
		Put("/api/v1/{project}/item/{itemId}")
	return &rs, err
}

//SaveLog attaches log in RP
func (c *Client) SaveLog(log *SaveLogRQ) (*EntryCreatedRS, error) {
	var rs EntryCreatedRS
	_, err := c.http.R().
		SetPathParams(map[string]string{
			"project": c.project,
		}).
		SetBody(log).
		SetResult(&rs).
		Post("/api/v1/{project}/log")
	return &rs, err
}

//GetLaunches retrieves latest launches
func (c *Client) GetLaunches() (*LaunchPage, error) {
	var launches LaunchPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

//GetLaunchesByFilter retrieves launches by filter
func (c *Client) GetLaunchesByFilter(filter map[string]string) (*LaunchPage, error) {
	var launches LaunchPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		SetQueryParams(filter).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

//GetLaunchesByFilterString retrieves launches by filter as string
func (c *Client) GetLaunchesByFilterString(filter string) (*LaunchPage, error) {
	var launches LaunchPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		SetQueryString(filter).
		Get("/api/v1/{project}/launch")
	return &launches, err
}

//GetLaunchesByFilterName retrieves launches by filter name
func (c *Client) GetLaunchesByFilterName(name string) (*LaunchPage, error) {
	filter, err := c.GetFiltersByName(name)
	if nil != err {
		return nil, err
	}

	if filter.Page.Size < 1 {
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

//GetFiltersByName retrieves filter by its name
func (c *Client) GetFiltersByName(name string) (*FilterPage, error) {
	var filter FilterPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project, "name": name}).
		SetQueryParam("filter.eq.name", name).
		SetResult(&filter).
		Get("/api/v1/{project}/filter")
	return &filter, err
}

//MergeLaunches merge two launches
func (c *Client) MergeLaunches(rq *MergeLaunchesRQ) (*LaunchResource, error) {
	var rs LaunchResource
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetBody(rq).
		SetResult(&rs).
		Post("/api/v1/{project}/launch/merge")
	return &rs, err
}
