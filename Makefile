sqlc:
	sqlc generate
.PHONY: sqlc

proto:
	protoc --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative internal/controller/protos/leaderboard.proto
.PHONY: proto

test:
	go test ./...
.PHONY: test
