FROM golang:1.23-alpine AS build

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o main cmd/api/main.go

FROM alpine:3.20.1 AS prod
WORKDIR /app

RUN apk add --no-cache tzdata

COPY --from=build /app/main /app/main

ENV TZ=America/Sao_Paulo

EXPOSE ${PORT}
CMD ["./main"]


