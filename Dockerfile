FROM golang:1.23.4

WORKDIR /smtp-server

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN mkdir app
COPY app ./app
COPY main.go ./
COPY authfile ./
COPY .env ./
COPY docker-build.sh ./
RUN chmod +x docker-build.sh
RUN ./docker-build.sh

RUN go build -o go-smtp-server

CMD ["./go-smtp-server"]

