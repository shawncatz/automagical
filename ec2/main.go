package ec2

import (
	"context"
	"errors"

	"github.com/sirupsen/logrus"
)

const (
	waitMax    = 50
	waitPoll   = 2
	tagName    = "automagical"
	tagVolume  = tagName + ":volume"
	tagRecord  = tagName + ":record"
	tagAddress = tagName + ":address"
)

func Handle(evt Event, ctx context.Context) (bool, error) {
	cfg := NewConfig()
	svc := NewService(evt.Region)
	h := &Handler{event: evt, ctx: ctx, config: cfg, service: svc}

	switch evt.Detail.State {
	case "running":
		err := h.Running()
		return err == nil, err
	case "terminated":
		err := h.Terminated()
		return err == nil, err
	case "":
		logrus.Error("Event.Detail.State is empty")
		return false, errors.New("Event.Detail.State is empty")
	default:
		logrus.Infof("Do not handle '%s' state", evt.Detail.State)
		return true, nil
	}
}
