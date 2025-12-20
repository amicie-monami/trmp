# Этап сборки
FROM golang:1.23-alpine AS builder

WORKDIR /app

RUN apk add --no-cache gcc musl-dev sqlite-dev

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=1 GOOS=linux go build -a -installsuffix cgo -o main ./cmd

FROM alpine:latest

RUN apk --no-cache add ca-certificates sqlite-libs

WORKDIR /root/

COPY --from=builder /app/main .

COPY --from=builder /app/*.sql ./

EXPOSE 6580

RUN rm writers.db

CMD ["./main"]