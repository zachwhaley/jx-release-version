NAME := new-release-version
ORG := zachwhaley
VERSION := 2.0.0

all: build test

.PHONY: test
test:
	go test -v .

.PHONY: build
build:
	go build -v .

.PHONY: clean
clean:
	-$(RM) $(NAME)
