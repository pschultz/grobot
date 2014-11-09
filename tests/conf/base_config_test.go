package conf

import (
	. "github.com/fgrosse/gobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	. "github.com/fgrosse/gobot/conf"
	yaml "gopkg.in/yaml.v2"
)

var _ = Describe("BaseConfigType", func() {
	configFileFixture := "../fixtures/gobot.yml"

	BeforeEach(func() {
		// reset the configuration
		Config = BaseConfigType{}
	})

	It("should load the configuration from yaml", func() {
		content := YAMLContent(`
			gomock:
    			mock-folder: tests/mocks
				mocks:
		        	- tests/fixtures/mock_source1.go
		        	- tests/fixtures/mock_source2.go
		`)
		var c BaseConfigType
		err := yaml.Unmarshal(content, &c)
		Expect(err).NotTo(HaveOccurred())
		Expect(c.Gomock).NotTo(BeNil())
		Expect(c.Gomock.MockFolder).To(Equal("tests/mocks"))
		Expect(c.Gomock.Mocks).To(HaveLen(2))
		Expect(c.Gomock.Mocks[0]).To(Equal("tests/fixtures/mock_source1.go"))
		Expect(c.Gomock.Mocks[1]).To(Equal("tests/fixtures/mock_source2.go"))
	})

	Describe("load the configuration from file", func() {
		It("should parse the yaml content", func() {
			err := LoadFromFile(configFileFixture)
			Expect(err).NotTo(HaveOccurred())
			Expect(Config.Gomock).NotTo(BeNil())
			Expect(Config.Gomock.MockFolder).To(Equal("tests/mocks"))
			Expect(Config.Gomock.Mocks).To(HaveLen(2))
			Expect(Config.Gomock.Mocks[0]).To(Equal("tests/fixtures/mock_source1.go"))
			Expect(Config.Gomock.Mocks[1]).To(Equal("tests/fixtures/mock_source2.go"))
		})

		Context("with gomock configuration", func() {
			PIt("should register gomock as vendor binary", func() {
				err := LoadFromFile(configFileFixture)
				Expect(err).NotTo(HaveOccurred())
				//	Expect(Config.VendorBins).To(HaveKey("mockgen"))
				//	Expect(Config.VendorBins["mockgen"]).To(Equal("code.google.com/p/gomock/mockgen"))
			})
		})
	})
})
