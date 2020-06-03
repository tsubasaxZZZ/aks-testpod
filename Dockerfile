ARG GO_VERSION=1.12.6
FROM golang:${GO_VERSION}-alpine AS build-stage
RUN apk add --no-cache git
WORKDIR /src
COPY ./go.mod ./
RUN go mod download
COPY . .
RUN go build -o /signal signal.go

FROM alpine:3.9
COPY --from=build-stage /signal /
EXPOSE 80
ENTRYPOINT ["/signal"]