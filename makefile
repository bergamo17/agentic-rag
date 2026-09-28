include .env
export

postgres:
	docker run --name postgres16 -p ${DB_PORT}:5432 -e POSTGRES_USER=${DB_USER} -e POSTGRES_PASSWORD=${DB_PASSWORD} -e POSTGRES_DB=${DB_DATABASE} -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=${DB_USER} --owner=root ${DB_DATABASE}

dropdb:
	docker exec -it postgres16 dropdb ${DB_DATABASE}

migrateup:
	migrate -path db/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_DATABASE}?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migrations -database "postgresql://${DB_USER}:${DB_PASSWORD}@localhost:5432/${DB_DATABASE}?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test: 
	go test -v -cover ./...

psql:
	docker exec -it agentic-ai-db psql -U agentic -d agentic_ai_prototype

server: 
	go run cmd/server/main.go

.PHONY: postgres createdb dropdb migrateup migratedown sqlc test psql server