install:
	 go get ./... && go install github.com/onsi/ginkgo/v2/ginkgo

run:
	go run main.go

test:
	go test ./...

db:
	brew services start mongodb-community