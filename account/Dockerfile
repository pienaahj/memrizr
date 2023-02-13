# Start from base image
FROM golang:alpine as builder

# Set the current working directory inside the container
WORKDIR /go/src/app

# Set the module env on
ENV GO111MODULE=on

# Install the reflex package
RUN go install github.com/cespare/reflex@latest

# Copy go mod and sum files
COPY go.mod .
# - this is only needed if there are dependancies
#COPY go.sum . 

# Download all dependencies
RUN go mod download

# Copy source from current directory to working director
COPY . .

# Build the application - called run
RUN go build -o ./run .

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /root/

# Copy executable from builder
COPY --from=builder /go/src/app/run .

EXPOSE 8080
CMD ["./run"]
