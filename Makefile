build:
	go run main.go

test:
	ginkgo ./...

db:
	brew services start mongodb-community