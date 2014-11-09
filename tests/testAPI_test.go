package tests

import (
	. "github.com/fgrosse/gobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Test API", func() {

	Describe("YAML helpers", func() {
		It("should transform a given string into byte and remove unnecessary indention", func() {
			content := YAMLContent(`

				mock-folder: tests/mocks
				mocks:
					- tests/fixtures/mock_source1.go
					- tests/fixtures/mock_source2.go
			`)
			Expect(content).To(Equal([]byte(
				"mock-folder: tests/mocks\n" +
					"mocks:\n" +
					"    - tests/fixtures/mock_source1.go\n" +
					"    - tests/fixtures/mock_source2.go\n",
			)))
		})
	})
})
