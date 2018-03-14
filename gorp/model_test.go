package gorp

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("UnixTime", func() {
	It("Should correctly deserialize unix time", func() {
		const jsonStr = `"1512114178671"`
		const expTime = "2017-12-01T07:42:59+00:00"

		var unitTime Timestamp
		err := json.Unmarshal([]byte(jsonStr), &unitTime)
		Expect(err).ShouldNot(HaveOccurred())

		unitTime = Timestamp{unitTime.Truncate(1 * time.Minute)}

		d, _ := time.Parse(time.RFC3339, expTime)
		d = d.In(time.Local).Truncate(1 * time.Minute)

		Expect(unitTime.Time).To(Equal(d))
	})
	It("Should return error on incorrect unix time", func() {
		const jsonStr = `"hello-world"`

		var unitTime Timestamp
		err := json.Unmarshal([]byte(jsonStr), &unitTime)
		Expect(err).Should(HaveOccurred())
	})
})

var _ = Describe("Helpers", func() {
	It("Direction should be converter correctly", func() {
		Expect(directionToStr(true)).To(Equal("ASC"))
		Expect(directionToStr(false)).To(Equal("DESC"))
	})

	It("Should convert filters correctly", func() {
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
		Expect(fp).To(Equal(map[string]string{
			"filter.cnt.name": "value",
			"filter.eq.desc":  "valuedesc",
			"page.sort":       "name,DESC",
		}))
	})
})
