include app.env

network:
	docker network create adrouter

postgres:
	docker run --network adrouter --name postgres17 -p 5432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres

redis:
	docker run --network adrouter --name redis7 -p 6379:6379 -d redis

createdb:
	docker exec -it postgres17 createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ${POSTGRES_DB}

dropdb:
	docker exec -it postgres17 dropdb ${POSTGRES_DB}

migrateup:
	migrate -path db/migration -database ${DB_SOURCE} -verbose up

migratedown:
	migrate -path db/migration -database ${DB_SOURCE} -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

composeup:
	docker compose up --attach api

composedown:
	docker compose down


.PHONY: network postgres redis createdb dropdb migrateup migratedown sqlc test server composeup composedown