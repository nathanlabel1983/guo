Default:
	go build -o bin/guo_server.exe cmd/server/main.go

Test:
	go test ./...

Run:
	go run cmd/server/main.go