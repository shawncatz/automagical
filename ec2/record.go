package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/sirupsen/logrus"
)

func (s *InstanceService) AttachRecord(id, tagName, tagValue string) error {
	if tagValue == "" {
		logrus.Infof("%s:%s:%s tagValue is empty", tagName, tagValue, id)
		return nil
	}

	return nil
}

type ChangeRecordInput struct {
	Zone   string
	Action string
	Type   string
	TTL    int64
	Name   string
	Value  string
}

func (s *InstanceService) ChangeRecord(change *ChangeRecordInput) error {
	input := &route53.ChangeResourceRecordSetsInput{
		HostedZoneId: aws.String(change.Zone),
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				&route53.Change{
					Action: aws.String(change.Action),
					ResourceRecordSet: &route53.ResourceRecordSet{
						Name: aws.String(change.Name),
						Type: aws.String(change.Type),
						TTL:  &change.TTL,
						ResourceRecords: []*route53.ResourceRecord{
							&route53.ResourceRecord{Value: aws.String(change.Value)},
						},
					},
				},
			},
		},
	}

	if _, err := s.r53.ChangeResourceRecordSets(input); err != nil {
		return err
	}

	return nil
}
