include app.env

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=${POSTGRES_USER} -e POSTGRES_PASSWORD=${POSTGRES_PASSWORD} -d postgres

createdb:
	docker exec -it postgres createdb --username=${POSTGRES_USER} --owner=${POSTGRES_USER} ad_router

dropdb:
	docker exec -it postgres dropdb ad_router

migrateup:
	migrate -path db/migration -database ${DB_URL} -verbose up

migratedown:
	migrate -path db/migration -database ${DB_URL} -verbose down

sqlc:
	sqlc generate

.PHONY: postgres createdb dropdb migrateup migratedown