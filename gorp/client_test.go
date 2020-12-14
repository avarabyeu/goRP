package gorp

import (
	"net/http"
	"net/http/httptest"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("ReportPortal Client", func() {
	It("Creates correctly", func() {
		client := NewClient("http://host.com", "prj", "uuid")

		Expect(client.project).To(Equal("prj"))
		Expect(client.http.HostURL).To(Equal("http://host.com"))
		Expect(client.http.Token).To(Equal("uuid"))
	})

	It("Handles wrong status codes as errors", func() {
		server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.WriteHeader(http.StatusBadRequest)
		}))
		defer server.Close()

		client := NewClient(server.URL, "prj", "uuid")

		_, err := client.GetLaunches()
		Expect(err).Should(HaveOccurred())
	})
})
