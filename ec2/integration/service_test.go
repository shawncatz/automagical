package integration

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/shawncatz/automagical/ec2"
)

var _ = Describe("Service", func() {
	var (
		service *ec2.InstanceService
		id      string
	)
	BeforeEach(func() {
		service = ec2.NewService("us-west-2")
		id = "i-0f0d055b47b25bc3b"
	})
	Context("GetInstance", func() {
		It("can get instance", func() {
			instance, err := service.GetInstance(id)
			Expect(err).NotTo(HaveOccurred())
			Expect(instance).NotTo(BeNil())
			Expect(*instance.InstanceId).To(Equal(id))
		})
	})
	Context("GetInstanceZone", func() {
		It("can get instance zone", func() {
			zone, err := service.GetInstanceZone(id)
			Expect(err).NotTo(HaveOccurred())
			Expect(zone).To(Equal("us-west-2a"))
		})
	})
	Context("GetTags", func() {
		It("can get tags", func() {
			instance, err := service.GetInstance(id)
			Expect(err).NotTo(HaveOccurred())
			tags := service.GetTags(instance.Tags)
			Expect(tags["Name"]).To(Equal("automagical-0"))
		})
	})
	Context("Address", func() {
		It("can find address", func() {
			address, err := service.FindAddress(id, "automagical:address", "automagical-address-0")
			Expect(err).NotTo(HaveOccurred())
			tags := service.GetTags(address.Tags)
			Expect(tags["Name"]).To(Equal("automagical-0"))
		})
	})
	Context("Volume", func() {
		It("can find volume", func() {
			address, err := service.FindVolume(id, "automagical:volume", "automagical-volume-0")
			Expect(err).NotTo(HaveOccurred())
			tags := service.GetTags(address.Tags)
			Expect(tags["Name"]).To(Equal("automagical-0"))
		})
	})
})
