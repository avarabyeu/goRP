package gorp

import (
	"fmt"
	"strconv"
)

//ConvertToFilterParams converts RP internal filter representation to query string
func ConvertToFilterParams(filter *FilterResource) map[string]string {
	params := map[string]string{}
	for _, f := range filter.Entities {
		params[fmt.Sprintf("filter.%s.%s", f.Condition, f.Field)] = f.Value
	}

	if filter.SelectionParams != nil {
		if filter.SelectionParams.PageNumber != 0 {
			params["page.page"] = strconv.Itoa(filter.SelectionParams.PageNumber)
		}
		if filter.SelectionParams.Orders != nil {
			for _, order := range filter.SelectionParams.Orders {
				params["page.sort"] = fmt.Sprintf("%s,%s", order.SortingColumn, directionToStr(order.Asc))
			}
		}

	}

	return params
}

func directionToStr(asc bool) string {
	if asc {
		return "ASC"
	}
	return "DESC"

}
