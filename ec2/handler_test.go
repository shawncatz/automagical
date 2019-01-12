package ec2_test

import (
	"context"
	"encoding/json"
	"io/ioutil"

	"github.com/aws/aws-sdk-go/aws"
	awsec2 "github.com/aws/aws-sdk-go/service/ec2"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/shawncatz/automagical/ec2"
	"github.com/shawncatz/automagical/ec2/ec2fakes"
)

var _ = Describe("Handler", func() {
	var (
		runningInstance *awsec2.Instance
		taggedInstance  *awsec2.Instance
		//taggedAddress   *awsec2.Address
		running ec2.Event
		ctx     context.Context
		svc     *ec2fakes.FakeService
		handler *ec2.Handler
	)
	BeforeEach(func() {
		running = loadEvent("running")
		ctx = context.Background()
		svc = &ec2fakes.FakeService{}
		handler = ec2.NewHandler(running, ctx, nil, svc)
		runningInstance = &awsec2.Instance{
			State: &awsec2.InstanceState{Name: aws.String("running")},
			Tags: []*awsec2.Tag{
				&awsec2.Tag{
					Key:   aws.String("test"),
					Value: aws.String("value"),
				},
			},
		}
		taggedInstance = &awsec2.Instance{
			State: &awsec2.InstanceState{Name: aws.String("running")},
			Tags: []*awsec2.Tag{
				&awsec2.Tag{
					Key:   aws.String("automagical:address"),
					Value: aws.String("tag-test-1"),
				},
			},
		}
		//taggedAddress = &awsec2.Address{
		//	AllocationId: aws.String("blarg"),
		//	InstanceId:   nil,
		//	Tags: []*awsec2.Tag{
		//		&awsec2.Tag{
		//			Key:   aws.String("automagical:address"),
		//			Value: aws.String("tag-test-1"),
		//		},
		//	},
		//}
	})
	Context("Running events", func() {
		It("handles a found instance", func() {
			svc.WaitReturns(runningInstance, map[string]string{}, nil)

			err := handler.Running()
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("handles a not found instance", func() {
			err := handler.Running()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("instance not found"))
		})
		It("handles a tagged instance and attaches address", func() {
			svc.WaitReturns(taggedInstance, map[string]string{"automagical:address": "tag-test-1"}, nil)
			svc.AttachAddressReturns(nil)

			err := handler.Running()
			Expect(err).To(BeNil())
			Expect(svc.WaitCallCount()).To(Equal(1))
			Expect(svc.AttachAddressCallCount()).To(Equal(1))
		})
	})
})

func loadEvent(name string) ec2.Event {
	evt := ec2.Event{}

	file, err := ioutil.ReadFile("fixtures/" + name + ".json")
	if err != nil {
		//t.Error("could not read json file: ", err)
		return evt
	}

	err = json.Unmarshal(file, &evt)
	if err != nil {
		//t.Error("could not unmarshal json: ", err)
		return evt
	}

	return evt
}
