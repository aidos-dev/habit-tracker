FROM golang:1.17.1-alpine3.14

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql


RUN apk add postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

# build go app
RUN go mod download
RUN go build -o habit-tracker ./cmd/main.go

CMD ["./habit-tracker"]