FROM golang:1.19-alpine

WORKDIR /var/www/app/go-telegram-bot-todolist

RUN go install github.com/cosmtrek/air@v1.40.4;

COPY go.mod go.sum ./

RUN go mod download

CMD ["air", "-c", ".air.toml"]

EXPOSE 7531