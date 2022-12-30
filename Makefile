install:
	 go get ./... && go install github.com/onsi/ginkgo/v2/ginkgo

build:
	go run main.go

test:
	ginkgo ./...

db:
	brew services start mongodb-community