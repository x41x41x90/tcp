build-client:
	go mod tidy
	go build -o client.exe cmd/client/main.go


build-server:
	go mod tidy
	go build -o server cmd/server/main.go