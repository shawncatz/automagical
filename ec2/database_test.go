package ec2_test

import (
	"github.com/aws/aws-sdk-go/aws"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"

	"github.com/shawncatz/automagical/ec2/ec2fakes"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Database", func() {
	var db *ec2fakes.FakeDatabase
	Context("Storing", func() {
		BeforeEach(func() {
			db = &ec2fakes.FakeDatabase{}
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
	})
})
