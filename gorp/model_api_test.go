package gorp

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestUnixTimeDeserialization(t *testing.T) {
	const jsonStr = `"1512114178671"`
	const expTime = "2017-12-01T07:42:59+00:00"

	var unitTime Timestamp
	err := json.Unmarshal([]byte(jsonStr), &unitTime)
	assert.NoError(t, err)

	unitTime = Timestamp{unitTime.Truncate(1 * time.Minute)}

	d, _ := time.Parse(time.RFC3339, expTime)
	d = d.In(time.Local).Truncate(1 * time.Minute)

	assert.Equal(t, d, unitTime.Time)
}

func TestUnixTimeSerialization(t *testing.T) {
	const jsonStr = `1512114179000`
	const expTime = "2017-12-01T07:42:59+00:00"

	d, _ := time.Parse(time.RFC3339, expTime)
	bytes, err := json.Marshal(&Timestamp{d})
	assert.NoError(t, err)
	assert.Equal(t, jsonStr, string(bytes))
}

func TestErrOnIncorrectTime(t *testing.T) {
	const jsonStr = `"hello-world"`

	var unitTime Timestamp
	err := json.Unmarshal([]byte(jsonStr), &unitTime)
	assert.Error(t, err)
}

func TestDirectionConverter(t *testing.T) {
	assert.Equal(t, "ASC", directionToStr(true))
	assert.Equal(t, "DESC", directionToStr(false))
}

func TestFiltersConverter(t *testing.T) {
	fp := ConvertToFilterParams(&FilterResource{
		Entities: []*FilterEntity{
			{
				Field:     "name",
				Condition: "cnt",
				Value:     "value",
			},
			{
				Field:     "desc",
				Condition: "eq",
				Value:     "valuedesc",
			},
		},
		SelectionParams: &FilterSelectionParam{
			Orders: []*FilterOrder{
				{
					Asc:           false,
					SortingColumn: "name",
				},
			},
		},
	})
	assert.Equal(t, map[string]string{
		"filter.cnt.name": "value",
		"filter.eq.desc":  "valuedesc",
		"page.sort":       "name,DESC",
	}, fp)
}
