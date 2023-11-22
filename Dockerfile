FROM golang:alpine as builder

WORKDIR /build

COPY go.mod .


RUN go mod download

COPY . .

RUN go build -o timerunner .


FROM alpine:latest as runner

WORKDIR /app

# copy binary and conf
COPY --from=builder /build/timerunner .
COPY conf.json .

# copy data dir
COPY data data


EXPOSE  3004


ENTRYPOINT ["/app/timerunner"]

