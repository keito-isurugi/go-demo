FROM golang:1.25

WORKDIR /app

RUN go install github.com/cosmtrek/air@v1.49.0

COPY . .

RUN go mod download

CMD ["air", "-c", ".air.toml"]
