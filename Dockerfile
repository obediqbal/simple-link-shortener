FROM golang:latest
COPY . /app
WORKDIR /app

RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq

WORKDIR /app/cmd

RUN go build -o simplelinkshortener

EXPOSE 8080

CMD ["./simplelinkshortener"]