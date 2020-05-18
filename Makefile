NAME ?= automagical
VERSION = $(shell cat VERSION)
FILE ?= $(NAME)-$(VERSION).zip
GITUSER ?= shawncatz
GITREPO ?= automagical

test: generate ginkgo

all: build

build: test $(NAME) $(FILE)

$(NAME):
	@GOOS=linux go build -o $(NAME)

$(FILE):
	zip $(FILE) ./$(NAME)

release: clean build
	git tag -f v$(VERSION)
	git push --tags

	github-release release \
        --user $(GITUSER) \
        --repo $(GITREPO) \
        --tag v$(VERSION) \
        --name "$(NAME)-v$(VERSION)"

	github-release upload \
        --user $(GITUSER) \
        --repo $(GITREPO) \
        --tag v$(VERSION) \
        --name "$(FILE)" \
        --file "$(FILE)"

release-delete:
	github-release delete -u $(GITUSER) -r $(GITREPO) -t v$(VERSION)

release-info:
	github-release info -u $(GITUSER) -r $(GITREPO)

clean:
	rm -f $(NAME) $(NAME)*.zip $(FILE)

generate:
	@go generate ./...

ginkgo:
	@ginkgo -nodes=5 ec2

integration:
	@aws-okta exec hub -- ginkgo -nodes=5 ec2/integration

deps:
	go get -u github.com/aktau/github-release
	go get -u github.com/onsi/ginkgo
	go get -u github.com/onsi/gomega
