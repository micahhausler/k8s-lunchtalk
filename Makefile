REPO=k8s-lunchtalk
SHA = $(shell git rev-parse --short HEAD)
GO_PKGS=$$(go list ./... | grep -v vendor)


.PHONY: setup fmt test test-cover vendored clean

all: test build

setup:
	go get golang.org/x/tools/cmd/cover
	go get -u github.com/kardianos/govendor
	go get -u github.com/golang/lint/golint

fmt:
	go fmt $(GO_PKGS)

build: fmt
	go build

test: fmt
	go test -race $(GO_PKGS)

test-cover: fmt
	go test -cover $(GO_PKGS)

docker:
	GOOS=linux go build
	docker build -t micahhausler/k8s-lunchtalk:$(SHA) .

push: docker
	docker tag micahhausler/k8s-lunchtalk:$(SHA) micahhausler/k8s-lunchtalk:latest
	docker push micahhausler/k8s-lunchtalk:$(SHA)
	docker push micahhausler/k8s-lunchtalk:latest

clean:
	rm ./$(REPO)
