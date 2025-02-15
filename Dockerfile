FROM golang:1.24.0-alpine as builder

WORKDIR /api

# COPY go.mod go.sum ./
COPY go.mod ./
RUN go mod download && go mod verify

COPY . ./

RUN CGO_ENABLED=0 GOOS=linux go build -o /app .

FROM alpine:latest

COPY --from=builder /app /app

CMD ["./app"]