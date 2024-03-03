all: .bin
	go get

.bin:
	go install github.com/vektra/mockery/v2@v2.40.0
	go install gotest.tools/gotestsum@v1.8.1

test:
	gotestsum -f dots -- -failfast -covermode=count -coverprofile coverage.out ./...
