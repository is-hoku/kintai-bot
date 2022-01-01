FROM golang:latest

RUN apt update && apt upgrade -y
RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.11.0/migrate.linux-amd64.tar.gz | tar xvz
RUN mv ./migrate.linux-amd64 /usr/bin/migrate

WORKDIR /go/src

COPY . .

RUN go mod download
RUN go get -u github.com/cosmtrek/air
RUN chmod 777 ./script/migrate-up
RUN chmod 777 ./script/migrate-down

WORKDIR ./app

CMD [ "air", "-c", ".air.toml" ]
