NAME := new-release-version
ORG := zachwhaley
VERSION := 2.1

all: build test

.PHONY: test
test:
	go test -v .

.PHONY: build
build:
	go build -v -ldflags '-X "main.Version=$(VERSION)-dev"' .

.PHONY: clean
clean:
	-$(RM) $(NAME)
