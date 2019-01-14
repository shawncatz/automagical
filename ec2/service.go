package ec2

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
)

//go:generate counterfeiter . Service
type Service interface {
	GetInstance(id string) (*ec2.Instance, error)
	GetInstanceZone(id string) (string, error)
	GetTags([]*ec2.Tag) map[string]string

	FindVolume(id, tagName, tagValue string) (*ec2.Volume, error)
	FindAddress(id, tagName, tagValue string) (*ec2.Address, error)

	AttachVolume(*ec2.Instance, *ec2.Volume) error
	AttachAddress(*ec2.Instance, *ec2.Address) error
	AttachRecord(id, tagName, tagValue string) error
}

type InstanceService struct {
	session *session.Session
	ec2     *ec2.EC2
	r53     *route53.Route53
	Region  string
}

func NewService(region string) *InstanceService {
	ses := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	return &InstanceService{
		session: ses,
		ec2:     ec2.New(ses, &aws.Config{Region: aws.String(region)}),
		r53:     route53.New(ses, &aws.Config{Region: aws.String(region)}),
		Region:  region,
	}
}
