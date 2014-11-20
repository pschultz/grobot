package tests

import (
	. "github.com/fgrosse/grobot/tests/mocks"
	. "github.com/fgrosse/grobot/tests/testAPI"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"code.google.com/p/gomock/gomock"
	"encoding/json"
	"fmt"
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
		Expect(conf.Version.String()).To(Equal("0.6"))
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

	It("should return an error if the minimum bot version from the configuration file is greater than the given bot version", func() {
		mockCtrl := gomock.NewController(GinkgoT())
		fileSystem := NewMockFileSystem(mockCtrl)
		grobot.FileSystemProvider = fileSystem
		configFilePath := "test-config.json"
		AssertFileWithContentExists(configFilePath, `{ "version": "1.23" }`, AnyTime, fileSystem)
		currentVersion := grobot.NewVersion("0.5")

		err := grobot.LoadConfigFromFile(configFilePath, currentVersion)
		Expect(err).To(HaveOccurred())
		Expect(err.Error()).To(Equal(fmt.Sprintf(`Error while read configuration file %s : The minimum required bot version is "1.23" but you are running bot version "0.5"`, configFilePath)))

		grobot.Reset()
		mockCtrl.Finish()
	})
})