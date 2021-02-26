FROM golang

ADD . /go/src/

WORKDIR /app

COPY . /app

RUN go build -o app client/*

ENTRYPOINT  /app/app