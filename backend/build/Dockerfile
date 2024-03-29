# Build stage
FROM golang:1.20.1-alpine3.17 AS build

LABEL build=temp

RUN go version

WORKDIR /app

COPY go.mod ./
COPY go.sum ./

# COPY ./ ./
# COPY migrate ./

COPY /backend ./backend
COPY /pkg ./pkg

RUN ls -lah


# build go app
RUN go mod download

# install migrate to do database migration
# RUN go install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest

RUN go build -o habit-tracker ./backend/cmd/habit/main.go




# Final stage
FROM alpine:3.17

RUN apk update

# install psql
RUN apk add --no-cache postgresql-client

WORKDIR /app

COPY --from=build /app/habit-tracker .
COPY backend/build/wait-for-postgres.sh /app/build/
COPY backend/build/.env /app/build/
COPY backend/configs/config.yml /app/configs/

RUN ls -lah
RUN ls -lah build
RUN ls -lah configs

# make wait-for-postgres.sh executable
RUN chmod +x /app/build/wait-for-postgres.sh


CMD ["./habit-tracker"]