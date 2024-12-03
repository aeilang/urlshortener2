postgresURL="postgres://lang:password@localhost:5432/urldb?sslmode=disable"

migrate_up:
	migrate -path="./database/migrate" -database=${postgresURL} up

migrate_down:
	migrate -path="./database/migrate" -database=${postgresURL} drop -f

lanch_postgres:
	docker run --name postgres-url \
	-e POSTGRES_USER=lang \
	-e POSTGRES_PASSWORD=password \
	-e POSTGRES_DB=urldb \
	-p 5432:5432 \
	-d postgres

lanch_redis:
	docker run --name redis \
	-p 6379:6379 \
	-d redis
