package integration

import (
	"context"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"github.com/shawncatz/automagical/ec2"
)

var _ = Describe("Handler", func() {
	var (
		evt     ec2.Event
		handler *ec2.Handler
	)
	BeforeEach(func() {
		evt = loadEvent("running")
		handler = ec2.NewHandler(evt, context.Background(), nil, nil, nil)
	})
	Context("Running", func() {
		It("it runs successfully", func() {
			err := handler.Running()
			Expect(err).NotTo(HaveOccurred())
		})
	})
})
