BINARY_PATH=_build/bin/resume-backend

.PHONY: clean fmt build
all: clean fmt build

.PHONY: fmt
fmt:
	go fmt ./.../

.PHONY: vet
vet:
	go vet ./.../

.PHONY: test
test:
	go test ./.../

.PHONY: build
build: clean test
	GOOS=linux GOARCH=amd64 go build -o $(BINARY_PATH) main.go

.PHONY: clean
clean:
	rm -f $(BINARY_PATH)
