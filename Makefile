format:
	go fmt ./...

compile:
	go build ./...

test:
	go test -v -cover ./...
