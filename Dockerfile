FROM golang:1.21-alpine as builder
WORKDIR /usr/local/src
RUN apk --no-cache add bash git build-base
COPY ["go.mod", "go.sum", "./"]
RUN go mod download
COPY . .
RUN CGO_ENABLED=1 go build -o /usr/local/src/bin/app ./cmd/SarkorTelecom-testTask/main.go

FROM alpine:latest as runner
WORKDIR /usr/local/src
COPY --from=builder /usr/local/src/bin/app .
COPY internal/config/.env internal/config/
CMD ["./app"]
