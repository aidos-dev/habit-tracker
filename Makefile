include backend/build/.env
export DB_PASSWORD

migrate:
	migrate -path ./migrations -database 'postgres://postgres:${DB_PASSWORD}@localhost:5432/postgres?sslmode=disable' up

migrate-down:
	migrate -path ./migrations -database 'postgres://postgres:${DB_PASSWORD}@localhost:5432/postgres?sslmode=disable' down

compose:
	docker-compose -f docker-compose.yml up --build habit-tracker
	docker image prune -f

recompose:
	docker rm build_habit-tracker_1
	docker rmi build_habit-tracker:latest
	docker image prune -f
	docker-compose -f build/docker-compose.yml up --build habit-tracker