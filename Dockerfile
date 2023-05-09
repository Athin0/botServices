FROM golang:latest

WORKDIR /botService
COPY go.* ./
RUN go mod download

COPY ./ /botService

RUN apt-get update && apt-get -y upgrade

RUN go build -o myapp ./cmd/bot.go

CMD ["/botService/myapp"]