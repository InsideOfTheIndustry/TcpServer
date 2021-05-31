FROM golang:1.16 AS build
WORKDIR /go/src/tcpserver
COPY . .
RUN go mod download
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o app main.go

FROM alpine:latest
WORKDIR /app
RUN mkdir ./config
COPY ./config/config.yaml ./config
COPY --from=build /go/src/tcpserver/app .

EXPOSE 4000 4001

CMD [ "./app" ]