all: .bin
	go get

.bin:
	go install github.com/vektra/mockery/v2@v2.40.0
