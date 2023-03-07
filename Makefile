vet:
	go vet ./...
.PHONY: vet

test:
	go test -v ./...
.PHONY: test

sync-vendor:
	go mod tidy && go mod vendor
.PHONY: sync-vendor