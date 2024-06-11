FROM golang:latest

WORKDIR /app

COPY ./ ./

RUN go mod download

RUN go build -o /postandcomments ./cmd/postandcomments

EXPOSE 8080 

CMD [ "/postandcomments", "--storage-type", "postgres"]