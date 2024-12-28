FROM golang:1.23.4

WORKDIR /smtp-server

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN mkdir app
COPY app ./app
RUN mkdir migrations
COPY migrations ./migrations
COPY main.go ./
COPY authfile ./
COPY .env ./

RUN go build -o go-smtp-server

CMD ["./go-smtp-server"]

