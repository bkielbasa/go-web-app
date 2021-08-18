FROM golang:1.16 AS builder

RUN mkdir /build/
WORKDIR /build/
COPY go.* .
RUN go mod download
COPY . .
RUN go build -o webapp

FROM alpine
WORKDIR /app
COPY --from=builder /build/webapp /
ENTRYPOINT ./webapp
