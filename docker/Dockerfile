FROM golang:1.20.11

WORKDIR /app

COPY . /app

RUN go mod download

RUN go build -o main

EXPOSE 8080

CMD ["./main"]

# TODO: use alpine to run ./main