FROM golang:1.24rc2-bookworm

RUN go version
ENV GOPATH=/

COPY ./ ./

# install psql
RUN apt-get update
RUN apt-get -y install postgresql-client

# make wait-for-postgres.sh executable
RUN chmod +x wait-for-postgres.sh

RUN go mod download
RUN go build -o pet-store ./cmd/api/main.go
RUN go get github.com/golang-migrate/migrate/v4

CMD ["./pet-store"]