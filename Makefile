include build/.env
export DB_PASSWORD

migrate:
	migrate -path ./migrations -database 'postgres://postgres:${DB_PASSWORD}@localhost:5432/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://postgres:${DB_PASSWORD}@localhost:5432/postgres?sslmode=disable' down

compose:
	docker-compose -f build/docker-compose.yml up --build habit-tracker

recompose:
	docker rm build_habit-tracker_1
	docker rmi build_habit-tracker:latest
	docker-compose -f build/docker-compose.yml up --build habit-tracker