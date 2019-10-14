.PHONY: build clean rebuild

pubsubroller: main.go create.go delete.go 
	GO111MODULES=on
	go build -o pubsubroller main.go delete.go create.go

build: pubsubroller

clean:
	rm -f pubsubroller
