package integration

import (
	"github.com/aws/aws-sdk-go/aws"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/shawncatz/automagical/ec2"
)

var _ = Describe("Database", func() {
	Context("Storing", func() {
		var (
			db    *ec2.InstanceDatabase
			table = "automagical_ec2"
		)
		BeforeEach(func() {
			db = ec2.NewDatabase(table)
		})
		It("Stores an instance", func() {
			err := db.Insert(&awsec2.Instance{
				InstanceId: aws.String("i-06db6eb9ed3ed4db5"),
				State:      &awsec2.InstanceState{Name: aws.String("running")},
				Tags: []*awsec2.Tag{
					&awsec2.Tag{
						Key:   aws.String("test"),
						Value: aws.String("value"),
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())
		})
		It("Retrieves an instance", func() {
			err := db.Insert(&awsec2.Instance{
				InstanceId: aws.String("i-06db6eb9ed3ed4db5"),
				State:      &awsec2.InstanceState{Name: aws.String("running")},
				Tags: []*awsec2.Tag{
					&awsec2.Tag{
						Key:   aws.String("test"),
						Value: aws.String("value"),
					},
				},
			})
			Expect(err).NotTo(HaveOccurred())

			instance, err := db.Find("i-06db6eb9ed3ed4db5")
			Expect(err).NotTo(HaveOccurred())
			Expect(instance).NotTo(BeNil())
			Expect(*instance.InstanceId).To(Equal("i-06db6eb9ed3ed4db5"))
		})
		It("Errors when retrieving an instance that doesn't exist", func() {
			_, err := db.Find("i-0123456789abcdef0")
			Expect(err).To(HaveOccurred())
			Expect(err.Error()).To(Equal("not found"))
		})
		It("Removes an instance", func() {
			err := db.Remove("i-06db6eb9ed3ed4db5")
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
