FROM golang:1.17

WORKDIR /go/src/app
COPY . .

COPY app/config.yml /go/src/app/config.yml

RUN go build -o /build/bin/api ./app

CMD ["/build/bin/api"]
