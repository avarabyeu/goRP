package gorp

import (
	"context"
	"fmt"

	"resty.dev/v3"

	"github.com/reportportal/goRP/v5/pkg/openapi"
)

const baseUrlLaunch = "/api/v1/{project}/launch"

type launchClient struct {
	http *resty.Client
}

// GetLaunchesByFilterString retrieves launches by filter as string
func (c *launchClient) GetLaunchesByFilterString(
	ctx context.Context,
	project, filter string,
) (*openapi.PageLaunchViewModel, error) {
	var launches openapi.PageLaunchViewModel
	_, err := c.http.R().SetContext(ctx).
		SetPathParam("project", project).
		SetResult(&launches).
		SetQueryString(filter).
		Get(baseUrlLaunch)
	return &launches, err
}

// GetAllLaunchesByFilterString retrieves every launch matching filter by iterating all pages.
// It starts at page 1 and keeps requesting the next page until the last page is reached.
// The filter string must not include page.page; it is managed internally.
func (c *launchClient) GetAllLaunchesByFilterString(
	ctx context.Context,
	project, filter string,
) ([]openapi.LaunchViewModel, error) {
	var all []openapi.LaunchViewModel
	for page := int64(1); ; page++ {
		pageFilter := fmt.Sprintf("page.page=%d&%s", page, filter)
		result, err := c.GetLaunchesByFilterString(ctx, project, pageFilter)
		if err != nil {
			return nil, err
		}
		all = append(all, result.Content...)
		if result.Page == nil || page >= result.Page.GetTotalPages() {
			break
		}
	}
	return all, nil
}

// GetLaunchesByFilterName retrieves launches by filter name
func (c *launchClient) GetLaunchesByFilterName(
	ctx context.Context,
	project, name string,
) (*openapi.PageLaunchViewModel, error) {
	filter, err := (&filterClient{http: c.http}).GetFiltersByName(ctx, project, name)
	if err != nil {
		return nil, err
	}

	if filter.Page.GetSize() < 1 || len(filter.Content) == 0 {
		return nil, fmt.Errorf(
			"no filter %s found",
			name,
		) //nolint:err113 //dynamic error is intentional
	}

	var launches openapi.PageLaunchViewModel
	params := ConvertToFilterParams(filter.Content[0])
	_, err = c.http.R().
		SetContext(ctx).
		SetPathParam("project", project).
		SetResult(&launches).
		SetQueryParamsFromValues(params).
		Get(baseUrlLaunch)
	return &launches, err
}
