GOARCH=amd64
CGO_ENABLED=0
GO_BUILD_FLAGS=-ldflags "-extldflags '-static'"

.PHONY: build
build: build-darwin build-linux build-windows

.PHONY: build-darwin
build-darwin:
	GOOS=darwin GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o bin/darwin/app ./cmd/app

.PHONY: build-linux
build-linux:
	GOOS=linux GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o bin/linux/app ./cmd/app

.PHONY: build-windows
build-windows:
	GOOS=windows GOARCH=$(GOARCH) CGO_ENABLED=$(CGO_ENABLED) go build $(GO_BUILD_FLAGS) -o bin/windows/app.exe ./cmd/app

.PHONY: clean
clean:
	rm -rf ./bin/app

.PHONY: format
format:
	go fmt $(go list ./... | grep -v /vendor/)

.PHONY: test
test:
	go vet $(go list ./... | grep -v /vendor/)
	go test -race $(go list ./... | grep -v /vendor/)


