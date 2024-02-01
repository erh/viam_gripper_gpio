
gripper-gpio: *.go cmd/module/*.go
	go build -o gripper-gpio cmd/module/cmd.go

test:
	go test

lint:
	gofmt -w -s .

module: gripper-gpio
	tar czf module.tar.gz gripper-gpio

all: module test
