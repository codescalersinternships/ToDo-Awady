FROM golang
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download && go mod verify

ENV DBFILE=./data/todo.db
ENV LISTENURL=:5000

COPY . .
RUN go build .
VOLUME [ "/data" ]

CMD ["./restapi"]