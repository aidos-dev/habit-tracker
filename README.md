# Habit tracker

### DataBase

In this project a PostgreSQL data base is applied. In order to start data base from docker container it is required to download a PostgreSQL docker image with the command:

```
docker pull postgres
```

To start the data base run the command:

```
docker run --name=habbit-db -e POSTGRES_PASSWORD='qwerty' -p 5432:5432 -d --rm postgres
```

To create migration files run the command:

```
migrate create -ext sql -dir ./schema -seq init
```

To make migration via the migration file run the command:

```
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up
```

To delete tables from the data base run the command:

```
migrate -path ./schema -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' down
```

To enter the data base run the command:

```
docker exec -it 0b3c8bef7b3d /bin/bash
```

Then inside the postgres docker container run the command:

```
psql -U postgres
```

Inside the postgres environment, to check all the tables run the command:

```
\d
```

To start the docker compose for the first time run the command:

```
docker-compose up --build habbit-tracker
```

When docker containers are built run the command without --build flag:

```
docker-compose up habbit-tracker
```
