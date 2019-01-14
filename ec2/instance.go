package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

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

func (s *InstanceService) GetTags(raw []*ec2.Tag) map[string]string {
	tags := make(map[string]string)

	for _, t := range raw {
		tags[*t.Key] = *t.Value
	}

	return tags
}
