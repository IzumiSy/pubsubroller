.PHONY: build clean rebuild

pubsubroller: main.go create_subscriptions.go create_topics.go delete_subscriptions.go delete_topics.go
	GO111MODULES=on
	go build -o pubsubroller \
		main.go \
		create_subscriptions.go \
		create_topics.go \
		delete_subscriptions.go \
		delete_topics.go

build: pubsubroller

clean:
	rm -f pubsubroller
