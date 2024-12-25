# Stage 1: Build the application
FROM golang:1.20-alpine AS build

# Install dependencies
RUN apk --no-cache add gcc g++ make ca-certificates

# Set the working directory
WORKDIR /go/src/github.com/akhilsharma90/go-graphql-microservice

# Copy Go module files
COPY go.mod go.sum ./

# Copy the vendor directory
COPY vendor vendor

# Copy application code
COPY account account
COPY catalog catalog
COPY order order

# Copy the SQL script into the build context
COPY up.sql order/up.sql

# Build the application
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./order/cmd/order

# Stage 2: Create the final lightweight image
FROM alpine:3.11

# Set the working directory
WORKDIR /usr/bin

# Copy the binary from the build stage
COPY --from=build /go/bin/app .

# Copy the SQL script into the final image
COPY --from=build /go/src/github.com/akhilsharma90/go-graphql-microservice/order/up.sql /usr/bin/up.sql

# Expose the application port
EXPOSE 8080

# Run the application
CMD ["app"]