FROM golang:1.18.0-buster

ARG USER=monke
ARG UID=1000
ARG GID=1000
# Default password for user
ARG PW=monke

WORKDIR /code

COPY go.mod ./

RUN go mod download && go mod verify

COPY . .

RUN go build -v -o monke /code ./

RUN useradd -m ${USER} --uid=${UID} && echo "${USER}:${PW}" | chpasswd

RUN chgrp -R ${GID} ./

USER ${USER} 

ENTRYPOINT ["./app.sh"]