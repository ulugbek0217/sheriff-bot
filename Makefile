postgres:
	docker run -d --name postgres16 -p 5555:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb --username=root --owner=root sheriff_bot
dropdb:
	docker exec -it postgres16 dropdb sheriff_bot

migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5555/sheriff_bot?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5555/sheriff_bot?sslmode=disable" -verbose down
migrateforce:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5555/sheriff_bot?sslmode=disable" force 1

sqlc:
	sqlc generate

prepare: postgres createdb migrateup sqlc

.PHONY: postgres createdb dropdb migrateup migratedown migrateforce sqlc