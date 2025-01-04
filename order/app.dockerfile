FROM golang:1.20-alpine AS build
RUN apk --no-cache add gcc g++ make ca-certificates
WORKDIR /go/src/github.com/akhilsharma90/go-graphql-microservice
COPY go.mod go.sum ./
COPY vendor vendor
COPY order order
RUN GO111MODULE=on go build -mod vendor -o /go/bin/app ./order/cmd/order

FROM alpine:3.11
WORKDIR /usr/bin
ENV ENV="dev"
COPY --from=build /go/bin/app .
EXPOSE 8080
CMD ["app"]