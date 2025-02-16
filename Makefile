
gripper-gpio: *.go cmd/module/*.go go.mod go.sum
	go build -tags netgo,osusergo -o gripper-gpio cmd/module/cmd.go

test:
	go test

lint:
	gofmt -w -s .

module.tar.gz: gripper-gpio
	tar czf module.tar.gz gripper-gpio

module: module.tar.gz

all: module test

update:
	go get go.viam.com/rdk@latest
	go mod tidy
