package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"github.com/fgrosse/grobot"
)

var _ = Describe("Configuration", func() {
	configData := []byte(`{
		"version": "0.6",
		"ginkgo": {
			"folder": "tests"
		}
	}`)

	It("should unmarshal the configuration version", func() {
		var conf grobot.Configuration
		err := json.Unmarshal(configData, &conf)
		Expect(err).NotTo(HaveOccurred())
		Expect(conf.Version).To(Equal("0.6"))
	})

	It("should unmarshal the module configurations", func() {
		var conf grobot.Configuration
		err := json.Unmarshal(configData, &conf)
		Expect(err).NotTo(HaveOccurred())
		Expect(conf.RawModuleConfigs).To(HaveLen(1))
		Expect(conf.RawModuleConfigs).To(HaveKey("ginkgo"))
	})

	It("should have a method to get raw configuration of a specific field", func() {
		var conf grobot.Configuration
		err := json.Unmarshal(configData, &conf)
		Expect(err).NotTo(HaveOccurred())

		rawConf, exists := conf.Get("ginkgo")
		Expect(exists).To(BeTrue())
		Expect(string(*rawConf)).To(Equal(string(`{
			"folder": "tests"
		}`)))

		_, exists = conf.Get("foobar")
		Expect(exists).To(BeFalse())
	})
})
