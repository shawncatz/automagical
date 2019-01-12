package ec2_test

import (
	"testing"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

func TestEc2(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Ec2 Suite")
}
