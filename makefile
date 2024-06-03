all:
	go build -v ./cmd/api
	migrate -path ./internal/db/migrations -database sqlite3://./internal/db/database.db up
