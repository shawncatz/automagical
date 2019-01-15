
all: build

build: generate test
	@GOOS=linux go build -o automagical

generate:
	@go generate ./...

test: ginkgo

ginkgo:
	@ginkgo -nodes=5 ec2

integration:
	@aws-okta exec hub -- ginkgo -nodes=5 ec2/integration
