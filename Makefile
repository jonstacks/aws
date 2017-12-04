
dep:
	dep ensure

install: dep
	go install ./cmd/...