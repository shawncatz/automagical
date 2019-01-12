package ec2

import (
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/sirupsen/logrus"
)

func (s *InstanceService) Wait(id string, max, poll time.Duration) (*ec2.Instance, error) {
	timeout := time.After(max * time.Second)
	tick := time.Tick(poll * time.Second)

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			logrus.Errorf("timed out waiting for instance %s", id)
			return nil, fmt.Errorf("timed out, Name tag not found for %s", id)
		// Got a tick, we should check
		case <-tick:
			ins, _ := s.GetInstance(id)
			if ins == nil {
				continue
			}

			// Wait for the running state, hopefully this means the tags are ready
			// in the past that wasn't always true
			// https://docs.aws.amazon.com/cli/latest/reference/ec2/wait/instance-running.html
			if *ins.State.Name != "running" {
				continue
			}

			return ins, nil
		}
	}
}

func (s *InstanceService) GetInstance(id string) (*ec2.Instance, error) {
	input := &ec2.DescribeInstancesInput{
		InstanceIds: []*string{aws.String(id)},
	}

	rsp, err := s.ec2.DescribeInstances(input)
	if err != nil {
		return nil, err
	}

	if len(rsp.Reservations) == 0 {
		return nil, fmt.Errorf("reservation for instance '%s' not found", id)
	}
	if len(rsp.Reservations[0].Instances) == 0 {
		return nil, fmt.Errorf("instance '%s' not found", id)
	}

	i := rsp.Reservations[0].Instances[0]

	return i, nil
}

func (s *InstanceService) GetInstanceZone(id string) (string, error) {
	input := &ec2.DescribeInstanceStatusInput{
		InstanceIds: []*string{aws.String(id)},
	}
	rsp, err := s.ec2.DescribeInstanceStatus(input)
	if err != nil {
		return "", err
	}
	if len(rsp.InstanceStatuses) != 1 {
		return "", nil // not everything cares about zone
	}

	return *rsp.InstanceStatuses[0].AvailabilityZone, nil
}

func getTags(raw []*ec2.Tag) map[string]string {
	tags := make(map[string]string)

	for _, t := range raw {
		tags[*t.Key] = *t.Value
	}

	return tags
}
