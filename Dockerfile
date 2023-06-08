FROM golang:1.20.1-alpine3.17

RUN go version

WORKDIR /app

COPY ./ ./

RUN ls -li

# install psql


RUN apk add postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x build/wait-for-postgres.sh



# build go app
RUN go mod download

# install migrate to do database migration
# RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go build -o habit-tracker ./cmd/habit/main.go

CMD ["./habit-tracker"]