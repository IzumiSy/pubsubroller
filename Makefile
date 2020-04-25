.PHONY: build clean test

pubsubroller: main.go
	 GO111MODULE=on go build -o pubsubroller

build: pubsubroller

clean:
	rm -f pubsubroller

test:
	GO111MODULE=on go test -v ./...
