package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (s *InstanceService) AttachAddress(instance *ec2.Instance, address *ec2.Address) error {
	input := &ec2.AssociateAddressInput{
		InstanceId:   instance.InstanceId,
		AllocationId: address.AllocationId,
	}

	if _, err := s.ec2.AssociateAddress(input); err != nil {
		return err
	}

	return nil
}

func (s *InstanceService) FindAddress(id, tagName, tagValue string) (*ec2.Address, error) {
	input := &ec2.DescribeAddressesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("tag:" + tagName),
				Values: []*string{aws.String(tagValue)},
			},
			&ec2.Filter{
				Name:   aws.String("domain"),
				Values: []*string{aws.String("vpc")},
			},
		},
	}

	out, err := s.ec2.DescribeAddresses(input)
	if err != nil {
		return nil, err
	}
	if len(out.Addresses) != 1 {
		return nil, fmt.Errorf("wrong number of addresses returned (%d) for %s:%s:%s", len(out.Addresses), id, tagName, tagValue)
	}
	if out.Addresses[0].InstanceId != nil && *out.Addresses[0].InstanceId != id {
		return nil, fmt.Errorf("address already attached to %s for %s:%s:%s", *out.Addresses[0].InstanceId, id, tagName, tagValue)
	}

	return out.Addresses[0], nil
}
