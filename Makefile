.PHONY: build clean test

pubsubroller: main.go \
	create.go \
	delete.go \
	config/config.go \
	google/adapters/subscription/subscription.go \
	google/adapters/topic/topic.go
	 GO111MODULE=on go build -o pubsubroller

build: pubsubroller

clean:
	rm -f pubsubroller

test:
	GO111MODULE=on go test -v ./...
