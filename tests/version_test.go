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
			Expect(version.GreaterThan(version)).To(BeFalse())
		})

		It("should be comparable to other versions with the same major and minor but different patch version", func() {
			lowerVersion := grobot.NewVersion("0.2.7")
			higherVersion := grobot.NewVersion("0.2.9")
			Expect(higherVersion.GreaterThan(lowerVersion)).To(BeTrue())
			Expect(lowerVersion.GreaterThan(higherVersion)).To(BeFalse())

			Expect(higherVersion.LowerThan(lowerVersion)).To(BeFalse())
			Expect(lowerVersion.LowerThan(higherVersion)).To(BeTrue())
		})

		It("should be comparable to other versions with the same major but different minor version", func() {
			lowerVersion := grobot.NewVersion("0.7.9")
			higherVersion := grobot.NewVersion("0.9.1")
			Expect(higherVersion.GreaterThan(lowerVersion)).To(BeTrue())
			Expect(lowerVersion.GreaterThan(higherVersion)).To(BeFalse())

			Expect(higherVersion.LowerThan(lowerVersion)).To(BeFalse())
			Expect(lowerVersion.LowerThan(higherVersion)).To(BeTrue())
		})

		It("should be comparable to other versions with the different major, minor and patch version", func() {
			lowerVersion := grobot.NewVersion("1.5.9")
			higherVersion := grobot.NewVersion("3.3.1")
			Expect(higherVersion.GreaterThan(lowerVersion)).To(BeTrue())
			Expect(lowerVersion.GreaterThan(higherVersion)).To(BeFalse())

			Expect(higherVersion.LowerThan(lowerVersion)).To(BeFalse())
			Expect(lowerVersion.LowerThan(higherVersion)).To(BeTrue())
		})

		It("should return an error when comparing different branch versions", func() {
			version1 := grobot.NewVersion("branch:foo")
			version2 := grobot.NewVersion("branch:bar")

			_, err := version1.GreaterThan(version2)
			Expect(err).To(HaveOccurred())

			_, err = version2.GreaterThan(version1)
			Expect(err).To(HaveOccurred())
		})

		It("should not return an error when comparing the same branch versions", func() {
			version1 := grobot.NewVersion("branch:foo")
			version2 := grobot.NewVersion("branch:foo")

			result, err := version1.GreaterThan(version2)
			Expect(err).NotTo(HaveOccurred())
			Expect(result).To(BeFalse())

			result, err = version1.LowerThan(version2)
			Expect(result).To(BeFalse())
		})
	})
})
