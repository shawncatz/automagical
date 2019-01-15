package ec2_test

import (
	"os"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/shawncatz/automagical/ec2"
)

var _ = Describe("Config", func() {
	Context("When no environment exists", func() {
		var c ec2.Config
		BeforeEach(func() {
			c = ec2.NewConfig()
		})
		It("uses the defaults", func() {
			Expect(c["table"]).To(Equal("automagical_ec2"))
		})
	})
	Context("When value is set in environment", func() {
		var c ec2.Config
		BeforeEach(func() {
			err := os.Setenv("AUTOMAGICAL_EC2_TABLE", "overridden_table_name")
			Expect(err).NotTo(HaveOccurred())
			c = ec2.NewConfig()
		})
		It("uses the override", func() {
			Expect(c["table"]).To(Equal("overridden_table_name"))
		})
	})
})
