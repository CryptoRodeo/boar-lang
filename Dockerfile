FROM golang:1.18.0-buster


WORKDIR /code

COPY go.mod ./

RUN go mod download && go mod verify

COPY . .

RUN chmod 777 ./app.sh

RUN go build -v -o /code ./

ENTRYPOINT ["./app.sh"]