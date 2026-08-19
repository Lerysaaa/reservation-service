FROM golang:1.26-alpine AS build

WORKDIR /app
COPY go.mod ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 go build -o /reservation-service ./cmd/api

FROM alpine:3.22
COPY --from=build /reservation-service /usr/local/bin/reservation-service
EXPOSE 8080
CMD ["reservation-service"]
