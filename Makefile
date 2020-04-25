.PHONY: build rebuild clean test update

pubsubroller: main.go
	 GO111MODULE=on go build -o pubsubroller

build: pubsubroller

rebuild: clean pubsubroller

clean:
	rm -f pubsubroller

update:
	GO111MODULE=on go get -u
	go mod tidy

test:
	GO111MODULE=on go test -v ./...
