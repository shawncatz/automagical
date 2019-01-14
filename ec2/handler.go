package ec2

import (
	"context"
	"fmt"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"

	"github.com/sirupsen/logrus"
)

type Handler struct {
	event   Event
	ctx     context.Context
	config  Config
	service Service
	Poll    time.Duration
	Max     time.Duration
}

func NewHandler(evt Event, ctx context.Context, cfg Config, svc Service) *Handler {
	if cfg == nil {
		cfg = NewConfig()
	}
	if svc == nil {
		svc = NewService(evt.Region)
	}
	return &Handler{
		event:   evt,
		ctx:     ctx,
		config:  cfg,
		service: svc,
		Poll:    waitPoll,
		Max:     waitMax,
	}
}

func (h *Handler) Running() error {
	id := h.event.Detail.Instance

	instance, err := h.Wait(id, waitMax, waitPoll)
	if err != nil {
		return err
	}

	tags := h.service.GetTags(instance.Tags)

	if err := h.Attach(instance, tags); err != nil {
		return err
	}

	return nil
}

func (h *Handler) Terminated() error {
	return nil
}

func (h *Handler) Attach(instance *ec2.Instance, tags map[string]string) error {
	errs := false

	if err := h.AttachAddress(instance, tagAddress, tags[tagAddress]); err != nil {
		logrus.Errorf("%s:%s:%s error %s", *instance.InstanceId, tagAddress, tags[tagAddress], err)
		errs = true
	}

	if err := h.AttachVolume(instance, tagVolume, tags[tagVolume]); err != nil {
		logrus.Errorf("%s:%s:%s error %s", *instance.InstanceId, tagVolume, tags[tagVolume], err)
		errs = true
	}

	if err := h.AttachRecord(instance, tagRecord, tags[tagRecord]); err != nil {
		logrus.Errorf("%s:%s:%s error %s", *instance.InstanceId, tagRecord, tags[tagRecord], err)
		errs = true
	}

	if errs {
		return fmt.Errorf("there were errors when processing attachments")
	}

	return nil
}

func (h *Handler) AttachAddress(instance *ec2.Instance, tagName, tagValue string) error {
	if tagValue == "" {
		logrus.Infof("%s:%s is empty", *instance.InstanceId, tagName)
		return nil
	}

	address, err := h.service.FindAddress(*instance.InstanceId, tagName, tagValue)
	if err != nil {
		return err
	}

	if err := h.service.AttachAddress(instance, address); err != nil {
		return err
	}

	return nil
}

func (h *Handler) AttachVolume(instance *ec2.Instance, tagName, tagValue string) error {
	if tagValue == "" {
		logrus.Infof("%s:%s is empty", *instance.InstanceId, tagName)
		return nil
	}

	volume, err := h.service.FindVolume(*instance.InstanceId, tagName, tagValue)
	if err != nil {
		return err
	}

	if err := h.service.AttachVolume(instance, volume); err != nil {
		return err
	}

	return nil
}

func (h *Handler) AttachRecord(instance *ec2.Instance, tagName, tagValue string) error {
	if tagValue == "" {
		logrus.Infof("%s:%s is empty", *instance.InstanceId, tagName)
		return nil
	}

	//volume, err := h.service.FindVolume(*instance.InstanceId, tagName, tagValue)
	//if err != nil {
	//	return err
	//}
	//
	//if err := h.service.AttachVolume(instance, volume); err != nil {
	//	return err
	//}

	return nil
}

func (h *Handler) Wait(id string, max, poll time.Duration) (*ec2.Instance, error) {

	ins, _ := h.service.GetInstance(id)
	if ins != nil && *ins.State.Name == "running" {
		return ins, nil
	}

	timeout := time.After(h.Max * time.Second)
	tick := time.Tick(h.Poll * time.Second)

	// Keep trying until we're timed out or got a result or got an error
	for {
		select {
		// Got a timeout! fail with a timeout error
		case <-timeout:
			logrus.Errorf("timed out waiting for instance %s", id)
			return nil, fmt.Errorf("timed out, running instance not found for %s", id)
		// Got a tick, we should check
		case <-tick:
			ins, _ := h.service.GetInstance(id)
			if ins == nil {
				continue
			}

			// Wait for the running state, hopefully this means the tags are ready
			// in the past that wasn't always true
			// https://docs.aws.amazon.com/cli/latest/reference/ec2/wait/instance-running.html
			if *ins.State.Name != "running" {
				continue
			}

			return ins, nil
		}
	}
}
