
.EXPORT_ALL_VARIABLES:
GO111MODULE = on

build:
	go build ./cmd/...

install:
	go install ./cmd/...

test:
	go test -v -cover -race ./...
