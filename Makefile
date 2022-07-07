createdb:
	createdb diplomski

dropdb:
	dropdb diplomski

sqlc:
	sqlc generate

migrateup:
	migrate -path db/migration --datebase @localhost 

.PHONY: createdb dropdb sqlc