include build/.env
export DB_PASSWORD

migrate:
	migrate -path ./migrations -database 'postgres://postgres:${DB_PASSWORD}@localhost:5432/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://postgres:${DB_PASSWORD}@localhost:5432/postgres?sslmode=disable' down