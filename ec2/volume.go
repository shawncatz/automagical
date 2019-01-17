package ec2

import (
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/ec2"
)

func (s *InstanceService) AttachVolume(instance *ec2.Instance, volume *ec2.Volume) error {
	input := &ec2.AttachVolumeInput{
		Device:     aws.String("/dev/sdf"), // TODO: support other devices and / or dynamically look up
		InstanceId: instance.InstanceId,
		VolumeId:   aws.String(*volume.VolumeId),
	}

	if _, err := s.ec2.AttachVolume(input); err != nil {
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
		return nil, fmt.Errorf("wrong number of volumes returned (%d) for %s:%s:%s, might already be attached", len(out.Volumes), id, tagName, tagValue)
	}
	if len(out.Volumes[0].Attachments) == 1 && *out.Volumes[0].Attachments[0].InstanceId != id {
		return nil, fmt.Errorf("volume already attached to %s for %s:%s:%s", *out.Volumes[0].Attachments[0].InstanceId, id, tagName, tagValue)
	}

	return out.Volumes[0], nil
}
