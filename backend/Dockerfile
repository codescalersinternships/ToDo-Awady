FROM golang AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

ENV DBFILE=./data/todo.db
ENV LISTENURL=:8080

COPY . .
RUN go build -a -ldflags "-linkmode external -extldflags '-static' -s -w" .


FROM alpine
WORKDIR /root/
COPY --from=builder /app/restapi ./
RUN mkdir data
VOLUME [ "/data" ]

CMD ["./restapi"]