# Habit tracker

<details>
<summary>Data Base</summary>
<br>

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
migrate create -ext sql -dir ./migrations -seq init
```

To make migration via the migration file run the command:

```
migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' up
```

To delete tables from the data base run the command:

```
migrate -path ./migrations -database 'postgres://postgres:qwerty@localhost:5432/postgres?sslmode=disable' down
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

In case of errors with migration and db becomes dirty, enter the data base and do the next:

```
select * from schema_migrations;
```

```
update schema_migrations set dirty =false where version=XXXX;
```

</details>

<details>
<summary>Docker-compose</summary>
<br>

To start the docker compose for the first time run the command:

```
docker-compose -f build/docker-compose.yml up --build habit-tracker
```

When docker containers are built run the command without --build flag:

```
docker-compose -f build/docker-compose.yml up habit-tracker
```

</details>

<details>
<summary>TG-bot</summary>
<br>

To run the telegram-bot service run the command:

```
./tg-bot-habit -tg-bot-token 'token'
```

</details>


<details>
<summary>Testing</summary>
<br>

The project is tested with unit testing and mocks

If the code logic is changed and test mocks and test funcs need to be updated accordingly, do the following:

    <details>
    <summary>Mock service</summary>
    <br>

    In order to generate the service layer mock do the following:

    1. Open the project in the terminal

    2. Go to directory of the service layer

    3. Run the command:

    ```
    go generate
    ```

    This command will regenerate the mock service

    4. Check if the unit tests needs to be updated accordning to updated logic

    </details>



</details>


