build:
	- chmod 777 scripts/protoc.sh && chmod 777 scripts/build.sh
	- scripts/build.sh

server:
	- go run cmd/server/main.go

client:
	- go run cmd/client/main.go

protoc:
	- scripts/protoc.sh

test:
	- go test internal/server/* -v -cover -race