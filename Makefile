run:
	go run main.go

test:
	go clean -testcache && go test -v -cover ./...

.PHONY: run test