install:
	 go get ./... && go install github.com/onsi/ginkgo/v2/ginkgo

build:
	go run main.go

test:
	go test ./...

db:
	brew services start mongodb-community