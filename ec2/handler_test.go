package ec2

import (
	"encoding/json"
	"io/ioutil"
	"testing"

	"github.com/aws/aws-sdk-go/service/ec2"
)

//func TestHandle(t *testing.T) {
//	evt := loadEvent("running", t)
//	ctx := context.Background()
//	res, err := Handle(evt, ctx)
//	if err != nil {
//		t.Errorf("Handle failed: %s", err)
//	}
//	if !res {
//		t.Error("Handle failed.")
//	}
//}

func TestHandler_Running(t *testing.T) {
	s := &InstanceService{
		client: &ec2.EC2{},
		region: "us-west-2",
	}
	h := &Handler{}
}

func TestHandler_Terminated(t *testing.T) {

}

func loadEvent(name string, t *testing.T) Event {
	evt := Event{}

	file, err := ioutil.ReadFile("fixtures/" + name + ".json")
	if err != nil {
		t.Error("could not read json file: ", err)
	}

	err = json.Unmarshal(file, &evt)
	if err != nil {
		t.Error("could not unmarshal json: ", err)
	}

	return evt
}

type ec2Mock struct {
}

func (e *ec2Mock) DescribeInstances(input *ec2.DescribeInstancesInput) (*ec2.DescribeInstancesOutput, error) {
	out := &ec2.DescribeInstancesOutput{
		Reservations: []*ec2.Reservation{
			&ec2.Reservation{},
		},
	}

	return out, nil
}
