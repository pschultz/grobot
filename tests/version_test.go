package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"github.com/fgrosse/grobot"
)

var _ = Describe("Version", func() {
	It("should unmarshable from JSON", func() {
		version := new(grobot.Version)
		data := []byte(`"1"`)
		err := json.Unmarshal(data, &version)
		Expect(err).NotTo(HaveOccurred())
		Expect(version.String()).To(Equal("1.0.0"))
		Expect(version.Major).To(Equal(1))
		Expect(version.Minor).To(Equal(0))
		Expect(version.Patch).To(Equal(0))

		version = new(grobot.Version)
		data = []byte(`"0.6"`)
		err = json.Unmarshal(data, &version)
		Expect(err).NotTo(HaveOccurred())
		Expect(version.String()).To(Equal("0.6.0"))
		Expect(version.Major).To(Equal(0))
		Expect(version.Minor).To(Equal(6))
		Expect(version.Patch).To(Equal(0))

		version = new(grobot.Version)
		data = []byte(`"0.6.1"`)
		err = json.Unmarshal(data, &version)
		Expect(err).NotTo(HaveOccurred())
		Expect(version.String()).To(Equal("0.6.1"))
		Expect(version.Major).To(Equal(0))
		Expect(version.Minor).To(Equal(6))
		Expect(version.Patch).To(Equal(1))
	})

	It("should marshal into the raw format", func() {
		version := grobot.NewVersion("0.6.1")
		data, err := json.Marshal(version)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal(`"0.6.1"`))

		version = grobot.NewVersion("0.6")
		data, err = json.Marshal(version)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal(`"0.6"`))

		version = grobot.NewVersion("0")
		data, err = json.Marshal(version)
		Expect(err).NotTo(HaveOccurred())
		Expect(string(data)).To(Equal(`"0"`))
	})

	It("should unmarshable from JSON when the `none` alias is used", func() {
		var version grobot.Version
		data := []byte(`"none"`)
		err := json.Unmarshal(data, &version)
		Expect(err).NotTo(HaveOccurred())
		Expect(version.String()).To(Equal("none"))
		Expect(version.Major).To(Equal(0))
		Expect(version.Minor).To(Equal(0))
		Expect(version.Patch).To(Equal(0))
	})

	It("should unmarshable from JSON when a branch is given", func() {
		var version grobot.Version
		data := []byte(`"branch:foo"`)
		err := json.Unmarshal(data, &version)
		Expect(err).NotTo(HaveOccurred())
		Expect(version.String()).To(Equal("branch:foo"))
		Expect(version.Major).To(Equal(0))
		Expect(version.Minor).To(Equal(0))
		Expect(version.Patch).To(Equal(0))
	})

	Context("comparing", func() {
		It("should never be greater than the same version", func() {
			version := grobot.NewVersion("0.7")
			Expect(version.GreaterThen(version)).To(BeFalse())
		})

		It("should be comparable to other versions with the same major and minor but different patch version", func() {
			lowerVersion := grobot.NewVersion("0.2.7")
			higherVersion := grobot.NewVersion("0.2.9")
			Expect(higherVersion.GreaterThen(lowerVersion)).To(BeTrue())
			Expect(lowerVersion.GreaterThen(higherVersion)).To(BeFalse())

			Expect(higherVersion.LowerThen(lowerVersion)).To(BeFalse())
			Expect(lowerVersion.LowerThen(higherVersion)).To(BeTrue())
		})

		It("should be comparable to other versions with the same major but different minor version", func() {
			lowerVersion := grobot.NewVersion("0.7.9")
			higherVersion := grobot.NewVersion("0.9.1")
			Expect(higherVersion.GreaterThen(lowerVersion)).To(BeTrue())
			Expect(lowerVersion.GreaterThen(higherVersion)).To(BeFalse())

			Expect(higherVersion.LowerThen(lowerVersion)).To(BeFalse())
			Expect(lowerVersion.LowerThen(higherVersion)).To(BeTrue())
		})

		It("should be comparable to other versions with the different major, minor and patch version", func() {
			lowerVersion := grobot.NewVersion("1.5.9")
			higherVersion := grobot.NewVersion("3.3.1")
			Expect(higherVersion.GreaterThen(lowerVersion)).To(BeTrue())
			Expect(lowerVersion.GreaterThen(higherVersion)).To(BeFalse())

			Expect(higherVersion.LowerThen(lowerVersion)).To(BeFalse())
			Expect(lowerVersion.LowerThen(higherVersion)).To(BeTrue())
		})

		// NOTE: this might not be final
		It("should return fals when comparing branch versions", func() {
			version1 := grobot.NewVersion("branch:foo")
			version2 := grobot.NewVersion("branch:bar")
			Expect(version1.GreaterThen(version2)).To(BeFalse())
			Expect(version2.GreaterThen(version1)).To(BeFalse())
			Expect(version1.LowerThen(version2)).To(BeFalse())
			Expect(version2.LowerThen(version1)).To(BeFalse())
		})
	})
})
