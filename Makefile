all: setup test

setup: .bin
	cp pre-commit .git/hooks/pre-commit
	go get

.bin:
	go install github.com/vektra/mockery/v2@v2.38.0
	go install gotest.tools/gotestsum@v1.8.1

test:
	gotestsum -f dots -- -failfast -covermode=count -coverprofile coverage.out ./...

	@go tool cover -func=coverage.out | grep 'total' | sed -e 's/\t\+/ /g' | sed -e 's/total: (statements)/TEST COVERAGE:/'
