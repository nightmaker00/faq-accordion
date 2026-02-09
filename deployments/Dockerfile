#building stage
FROM golang:1.22-alpine AS builder

ENV CGO_ENABLED=0 \
    GOOS=linux \
    GOARCH=amd64
    
WORKDIR /app
    
COPY go.mod go.sum ./
RUN go mod download
    

COPY . .
    

RUN go build -o app ./cmd/app
    
#runtime stage    
FROM alpine:3.19
    

RUN apk add --no-cache ca-certificates
    
WORKDIR /app

COPY --from=builder /app/app /app/app
    

EXPOSE 8080
    
ENTRYPOINT ["/app/app"]    