package main

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	"testing"
)

func TestCli(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "CLI Suite")
}

var _ = Describe("YES answer", func() {
	It("Should understand yes", func() {
		Expect(answerYes("yes")).To(Equal(true))
	})
	It("Should understand uppercase yes", func() {
		Expect(answerYes("YES")).To(Equal(true))
	})

	It("Should be false if answer is empty", func() {
		Expect(answerYes("")).To(Equal(false))
	})
	It("Should be false if answer is no", func() {
		Expect(answerYes("no")).To(Equal(false))
	})
})
