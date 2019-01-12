package ec2

import (
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/aws/aws-sdk-go/service/route53"
)

//go:generate counterfeiter . Service
type Service interface {
	GetTags([]*ec2.Tag) map[string]string
	GetInstance(id string) (*ec2.Instance, error)
	GetInstanceZone(id string) (string, error)
	Wait(id string, max, poll time.Duration) (*ec2.Instance, map[string]string, error)
	FindVolume(id, tagName, tagValue string) (*ec2.Volume, error)
	AttachVolume(id, tagName, tagValue string) error
	FindAddress(id, tagName, tagValue string) (*ec2.Address, error)
	AttachAddress(id, tagName, tagValue string) error
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
