FROM golang:1.24-alpine AS builder
WORKDIR /src
COPY go.mod ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/emergency-his .

FROM alpine:3.21
RUN apk add --no-cache ca-certificates tzdata
ENV TZ=Asia/Shanghai
WORKDIR /app
COPY --from=builder /out/emergency-his /usr/local/bin/emergency-his
EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/emergency-his"]
