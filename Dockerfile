FROM golang:1.23.4

ENV TZ=Europe/Rome

WORKDIR /smtp-server

COPY go.mod ./
COPY go.sum ./
RUN go mod download

RUN mkdir app
COPY app ./app
RUN mkdir ssl
COPY ssl ./ssl
RUN mkdir migrations
COPY migrations ./migrations
COPY main.go ./
COPY authfile ./
COPY .env ./

RUN apt-get update && apt-get install -y tzdata && rm -rf /var/lib/apt/lists/*

RUN ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && echo $TZ > /etc/timezone

RUN go build -o go-smtp-server

CMD ["./go-smtp-server"]

