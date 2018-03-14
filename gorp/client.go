package gorp

import (
	"fmt"
	"gopkg.in/resty.v1"
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
		SetHostURL(host).
		SetAuthToken(uuid).
		OnAfterResponse(func(client *resty.Client, rs *resty.Response) error {
		if (rs.StatusCode() / 100) >= 4 {
			return fmt.Errorf("status code error: %d", rs.StatusCode())
		}
		return nil
	})
	return &Client{
		project: project,
		http:    http,
	}
}

//GetLaunches retrieves latest launches
func (c *Client) GetLaunches() (LaunchPage, error) {
	var launches LaunchPage
	_, err := c.http.R().
		SetPathParams(map[string]string{"project": c.project}).
		SetResult(&launches).
		Get("/api/v1/{project}/launch")
	return launches, err
}
