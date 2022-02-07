FROM golang:1.17 as builder

WORKDIR /go/src

COPY . .

RUN go build -o /build/bin/api ./app

FROM centos:7

COPY --from=builder /build/bin/api /build/bin/api

COPY app/config.yml /go/src/config.yml

CMD ["/build/bin/api"]
