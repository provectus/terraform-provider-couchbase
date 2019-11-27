FROM golang:latest
ARG GOOS
ENV GOOS ${GOOS:-darwin}
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o /app/terraform-provider-couchbase .
