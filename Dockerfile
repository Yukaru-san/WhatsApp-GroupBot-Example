FROM golang:1.13-alpine as builder1

WORKDIR /go/src/main

# install required package(s)
RUN apk --no-cache add ca-certificates git

# Copy files
COPY . .

# Copy dependency list
RUN go get -d -v ./...

# Compile
RUN go build -o main

# Create new stage based on alpine
FROM alpine:latest

#Copy ca certs
COPY --from=builder1 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
# Copy compiled binary from builder1
WORKDIR /app
RUN mkdir /app/data/

COPY ./hangmanWordlists ./hangmanWordlists
COPY ./stickers ./stickers
COPY --from=builder1 /go/src/main/main .

# Set Debuglevel and start the server
ENV S_LOG_LEVEL debug
CMD [ "/app/main" ]
