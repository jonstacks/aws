
.EXPORT_ALL_VARIABLES:
GO111MODULE = on

install:
	go install ./cmd/...

test:
	go test -v -cover ./...
