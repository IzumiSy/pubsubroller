.PHONY: build clean

build: main.go create_subscriptions.go create_topics.go
	go build -o pubsubroller main.go create_subscriptions.go create_topics.go

clean: pubsubroller
	rm -r pubsubroller
