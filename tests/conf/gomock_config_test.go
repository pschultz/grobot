package conf

import (
	. "github.com/fgrosse/gobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/fgrosse/gobot/conf"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("GomockConfigType", func() {

	It("should load the configuration from yaml", func() {
		content := YAMLContent(`
			mock-folder: tests/mocks
			mocks:
				- tests/fixtures/mock_source1.go
				- tests/fixtures/mock_source2.go
		`)
		var c GomockConfigType
		err := yaml.Unmarshal(content, &c)
		Expect(err).NotTo(HaveOccurred())
		Expect(c.MockFolder).To(Equal("tests/mocks"))
		Expect(c.Mocks).To(HaveLen(2))
		Expect(c.Mocks[0]).To(Equal("tests/fixtures/mock_source1.go"))
		Expect(c.Mocks[1]).To(Equal("tests/fixtures/mock_source2.go"))
	})

})
