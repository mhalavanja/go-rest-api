postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:alpine

startDb:
	docker start postgres

createdb:
	docker exec -it postgres createdb --username=root --owner=root diplomski

dropdb:
	docker exec -it postgres dropdb diplomski

shell:
	docker exec -it postgres psql -d diplomski

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/diplomski?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/diplomski?sslmode=disable" -verbose down

sqlc:
	sqlc generate

server:
	go run main.go
.PHONY: postgres startDb createdb dropdb shell migrateup migratedown sqlc server
