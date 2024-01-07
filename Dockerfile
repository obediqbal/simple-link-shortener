FROM golang:latest
COPY . /app
WORKDIR /app

RUN go get github.com/gorilla/mux
RUN go get github.com/lib/pq
RUN go get golang.org/x/net/publicsuffix

WORKDIR /app/cmd

RUN go build -o simplelinkshortener

EXPOSE 8080

CMD ["./simplelinkshortener"]