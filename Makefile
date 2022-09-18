build:
	go build -o bin/transactionhistory main.go

run:
	go run main.go

test:
	go test ./...

clean:
	rm -rf bin

all:  test run

