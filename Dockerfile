FROM golang:alpine
RUN apk add --no-cache --upgrade bash build-base
WORKDIR /app
COPY ./scripts/wait-for.sh .
RUN chmod +x wait-for.sh
COPY server ./server
COPY go.mod go.sum configuration.json .golangci.yml main.go ./
RUN go get ./...
CMD ["go", "test", "./...", "-cover"]