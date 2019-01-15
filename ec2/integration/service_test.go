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
		id = "i-0665bb31b8d6ad04a"
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
			Expect(tags["Name"]).To(Equal("hub-openvpn"))
		})
	})
})
