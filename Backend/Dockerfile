# Build stage
FROM golang:1.23-alpine AS build

WORKDIR /app

RUN apk add --no-cache git openssh build-base pkgconfig
RUN mkdir /root/.ssh/ && \
  ssh-keyscan -t rsa github.com >> ~/.ssh/known_hosts

COPY go.mod .
COPY go.sum .

RUN --mount=type=ssh go mod download
 
COPY ./ ./
RUN go build -o ./tmp/main ./main.go


# Development stage

FROM build AS dev

RUN apk add --no-cache git
RUN go install github.com/air-verse/air@latest
 
EXPOSE 3000

WORKDIR /app

CMD ["air", "-c", ".air.toml"]