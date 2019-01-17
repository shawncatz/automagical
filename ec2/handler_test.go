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
		taggedAddress   *awsec2.Address
		running         ec2.Event
		svc             *ec2fakes.FakeService
		db              *ec2fakes.FakeDatabase
		handler         *ec2.Handler
		id              string
	)
	BeforeEach(func() {
		id = "i-0123456789abcdef0"
		running = loadEvent("running", id)
		svc = &ec2fakes.FakeService{}
		db = &ec2fakes.FakeDatabase{}
		handler = ec2.NewHandler(running, context.Background(), nil, svc, db)
		handler.Poll = 1
		handler.Max = 3
		runningInstance = &awsec2.Instance{
			InstanceId: aws.String(id),
			State:      &awsec2.InstanceState{Name: aws.String("running")},
			Tags: []*awsec2.Tag{
				&awsec2.Tag{
					Key:   aws.String("test"),
					Value: aws.String("value"),
				},
				&awsec2.Tag{
					Key:   aws.String("automagical"),
					Value: aws.String("true"),
				},
			},
		}
		taggedInstance = &awsec2.Instance{
			InstanceId: aws.String(id),
			State:      &awsec2.InstanceState{Name: aws.String("running")},
			Tags: []*awsec2.Tag{
				&awsec2.Tag{
					Key:   aws.String("automagical:address"),
					Value: aws.String("tag-test-1"),
				},
			},
		}
		taggedAddress = &awsec2.Address{
			AllocationId: aws.String("blarg"),
			InstanceId:   nil,
			Tags: []*awsec2.Tag{
				&awsec2.Tag{
					Key:   aws.String("automagical:address"),
					Value: aws.String("tag-test-1"),
				},
			},
		}
	})
	Context("Running events", func() {
		It("handles a found instance", func() {
			svc.GetInstanceReturns(runningInstance, nil)
			svc.GetTagsReturns(map[string]string{"automagical": "true"})

			err := handler.Running()
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("handles waiting for an instance", func() {
			svc.GetInstanceReturns(runningInstance, nil)
			svc.GetTagsReturnsOnCall(0, nil)
			svc.GetTagsReturnsOnCall(1, nil)
			svc.GetTagsReturnsOnCall(2, map[string]string{"automagical": "true"})

			err := handler.Running()
			Expect(err).ShouldNot(HaveOccurred())
		})
		It("handles a not found instance", func() {
			err := handler.Running()
			Expect(err).Should(HaveOccurred())
			Expect(err.Error()).To(Equal("timed out, running instance not found for " + id))
		})
		It("handles a tagged instance and attaches address", func() {
			svc.GetInstanceReturns(taggedInstance, nil)
			svc.GetTagsReturns(map[string]string{"automagical:address": "tag-test-1", "automagical": "true"})
			svc.FindAddressReturns(taggedAddress, nil)
			svc.AttachAddressReturns(nil)

			err := handler.Running()
			Expect(err).To(BeNil())
			Expect(svc.AttachAddressCallCount()).To(Equal(1))
			Expect(svc.GetInstanceCallCount()).To(Equal(1))
		})
	})
})

func loadEvent(name string, id string) ec2.Event {
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

	evt.Detail.Instance = id
	return evt
}
