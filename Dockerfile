# workspace (GOPATH) configured at /go
FROM golang:1.22.1-alpine3.19 AS builder

RUN mkdir -p /app
WORKDIR /app


COPY .env ./
COPY go.mod ./
COPY ./templates ./templates
COPY ./static ./static

# Copy the local package files to the container's workspace.
COPY . ./

# installing depends and build
RUN export CGO_ENABLED=0
RUN export GOOS=linux
RUN go build -o ./ /app/main.go
RUN mv main /

FROM alpine

RUN mkdir -p /app
WORKDIR /app

COPY --from=builder /app/.env /
COPY --from=builder /main /
COPY --from=builder /app/templates ./templates
COPY --from=builder /app/static ./static

EXPOSE 8080

CMD ["/main"]