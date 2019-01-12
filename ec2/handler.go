package ec2

import (
	"context"
	"errors"
	"fmt"

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

type Handler struct {
	event   Event
	ctx     context.Context
	config  Config
	service *InstanceService
}

func (h *Handler) Running() error {
	id := h.event.Detail.Instance
	instance, err := h.service.Wait(id, waitMax, waitPoll)
	if err != nil {
		return err
	}

	tags := getTags(instance.Tags)

	errs := false

	logrus.Infof("%s:%s", tagAddress, h.event.Detail.Instance)
	if err := h.service.AttachAddress(id, tagAddress, tags[tagAddress]); err != nil {
		logrus.Errorf("%s:%s error %s", tagAddress, h.event.Detail.Instance, err)
		errs = true
	}

	logrus.Infof("%s:%s", tagVolume, h.event.Detail.Instance)
	if err := h.service.AttachVolume(id, tagVolume, tags[tagVolume]); err != nil {
		logrus.Errorf("%s:%s error %s", tagVolume, h.event.Detail.Instance, err)
		errs = true
	}

	logrus.Infof("%s:%s", tagRecord, h.event.Detail.Instance)
	if err := h.service.AttachRecord(id, tagRecord, tags[tagRecord]); err != nil {
		logrus.Errorf("%s:%s error %s", tagRecord, h.event.Detail.Instance, err)
		errs = true
	}

	if errs {
		return fmt.Errorf("there were errors when processing attachments")
	}

	return nil
}

func (h *Handler) Terminated() error {
	return nil
}
