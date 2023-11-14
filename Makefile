.PHONY: all
all: vet build

.PHONY: build
build:
	go build ./cmd/qube

.PHONY: vet
vet:
	go vet ./...

.PHONY: test
test:
	go test -v -count=1 ./...

.PHONY: testacc
testacc:
	$(MAKE) test TEST_ACC=1

.PHONY: lint
lint:
	golangci-lint run
