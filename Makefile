.PHONY: build clean rebuild

pubsubroller: main.go create_subscriptions.go create_topics.go delete_subscriptions.go delete_topics.go
	go build -o pubsubroller main.go create_subscriptions.go create_topics.go delete_subscriptions.go delete_topics.go

build: pubsubroller

clean:
	rm -f pubsubroller
