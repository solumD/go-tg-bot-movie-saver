FROM golang:1.22

WORKDIR /usr/local/src/application

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ./ ./
RUN go build -o ./bin/app cmd/tg_bot/app.go

CMD ["./bin/app"]