package gorp

import (
	"encoding/json"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"time"
)

var _ = Describe("UnixTime", func() {
	It("Should correctly deserialize unit time", func() {
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
})
