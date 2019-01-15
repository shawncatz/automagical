
all: build

build: generate test
	@GOOS=linux go build -o automagical

generate:
	@go generate ./...

test:
	@go test -v ./...

ginkgo:
	@ginkgo -cover -nodes=5 ec2

integration:
	@ginkgo -nodes=5 ec2/integration
