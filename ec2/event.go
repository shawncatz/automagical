package ec2

type Event struct {
	Version    string
	ID         string
	DetailType string `json:"detail-type"`
	Source     string
	Account    string
	Time       string
	Region     string
	Test       bool
	Detail     EventDetail
	Resources  []string
}

type EventDetail struct {
	Instance string `json:"instance-id"`
	State    string
}
