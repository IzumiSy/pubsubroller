.PHONY: build clean test

pubsubroller: main.go create.go delete.go config/config.go subscription/subscription.go topic/topic.go
	go build -o pubsubroller

build: pubsubroller

clean:
	rm -f pubsubroller

test:
	go test
