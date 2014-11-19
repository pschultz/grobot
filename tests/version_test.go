package tests

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"encoding/json"
	"github.com/fgrosse/grobot"
)

var _ = Describe("Version", func() {
	It("should unmarshable from JSON", func() {
		var version grobot.Version
		data := []byte(`"0.6"`)
		err := json.Unmarshal(data, &version)
		Expect(err).NotTo(HaveOccurred())
		Expect(version.String()).To(Equal("0.6"))
		Expect(version.Major).To(Equal(0))
		Expect(version.Minor).To(Equal(6))
	})

	It("should unmarshable from JSON when the `none` alias is used", func() {
		var version grobot.Version
		data := []byte(`"none"`)
		err := json.Unmarshal(data, &version)
		Expect(err).NotTo(HaveOccurred())
		Expect(version.String()).To(Equal("none"))
		Expect(version.Major).To(Equal(0))
		Expect(version.Minor).To(Equal(0))
	})

	Context("comparing", func() {
		It("should never be greater than the same version", func() {
			version := grobot.NewVersion("0.7")
			Expect(version.GreaterThen(version)).To(BeFalse())
		})

		It("should be comparable to other versions with the same major but different minor version", func() {
			lowerVersion := grobot.NewVersion("0.7")
			higherVersion := grobot.NewVersion("0.9")
			Expect(higherVersion.GreaterThen(lowerVersion)).To(BeTrue())
			Expect(lowerVersion.GreaterThen(higherVersion)).To(BeFalse())

			Expect(higherVersion.LowerThen(lowerVersion)).To(BeFalse())
			Expect(lowerVersion.LowerThen(higherVersion)).To(BeTrue())
		})

		It("should be comparable to other versions with the differnet major but same minor version", func() {
			lowerVersion := grobot.NewVersion("1.2")
			higherVersion := grobot.NewVersion("2.2")
			Expect(higherVersion.GreaterThen(lowerVersion)).To(BeTrue())
			Expect(lowerVersion.GreaterThen(higherVersion)).To(BeFalse())

			Expect(higherVersion.LowerThen(lowerVersion)).To(BeFalse())
			Expect(lowerVersion.LowerThen(higherVersion)).To(BeTrue())
		})

		It("should be comparable to other versions with the different major and minor version", func() {
			lowerVersion := grobot.NewVersion("1.9")
			higherVersion := grobot.NewVersion("3.1")
			Expect(higherVersion.GreaterThen(lowerVersion)).To(BeTrue())
			Expect(lowerVersion.GreaterThen(higherVersion)).To(BeFalse())

			Expect(higherVersion.LowerThen(lowerVersion)).To(BeFalse())
			Expect(lowerVersion.LowerThen(higherVersion)).To(BeTrue())
		})
	})
})
