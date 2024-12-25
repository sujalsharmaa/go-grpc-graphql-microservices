# Stage 1: Build the application
FROM golang:1.21-alpine AS build

# Install dependencies
RUN apk --no-cache add gcc g++ make ca-certificates

# Set the working directory
WORKDIR /go/src/github.com/akhilsharma90/go-graphql-microservice

# Copy Go module files
COPY . .

# Build the application
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./account/cmd/account

# Stage 2: Create the final lightweight image
FROM alpine:3.11

# Set the working directory
WORKDIR /usr/bin

# Copy the binary from the build stage
COPY --from=build /go/bin/app .


# Expose the application port
EXPOSE 8080

# Run the application
CMD ["app"]