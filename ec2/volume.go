package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (s *InstanceService) AttachVolume(id, tagName, tagValue string) error {
	if tagValue == "" {
		return nil
	}

	volume, err := s.FindVolume(id, tagName, tagValue)
	if err != nil {
		return err
	}

	input := &ec2.AttachVolumeInput{
		Device:     aws.String("/dev/sdf"),
		InstanceId: aws.String(id),
		VolumeId:   aws.String(*volume.VolumeId),
	}

	_, err = s.ec2.AttachVolume(input)
	if err != nil {
		return err
	}

	return nil
}

func (s *InstanceService) FindVolume(id, tagName, tagValue string) (*ec2.Volume, error) {
	zone, err := s.GetInstanceZone(id)
	if err != nil {
		return nil, err
	}

	input := &ec2.DescribeVolumesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("tag:" + tagName),
				Values: []*string{aws.String(tagValue)},
			},
			&ec2.Filter{
				Name:   aws.String("status"),
				Values: []*string{aws.String("available")},
			},
			&ec2.Filter{
				Name:   aws.String("availability-zone"),
				Values: []*string{aws.String(zone)},
			},
		},
	}

	out, err := s.ec2.DescribeVolumes(input)
	if err != nil {
		return nil, err
	}
	if len(out.Volumes) != 1 {
		return nil, fmt.Errorf("wrong number of volumes returned (%d) for %s:%s:%s", len(out.Volumes), id, tagName, tagValue)
	}

	return out.Volumes[0], nil
}
