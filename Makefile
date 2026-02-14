.PHONY: all
all: vet test build

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

.PHONY: bench
bench:
	go test -run='0^' --count=3 -bench . -benchmem

.PHONY: lint
lint:
	golangci-lint run

.PHONY: demo
# see https://github.com/charmbracelet/vhs
demo:
	vhs demo.tape
