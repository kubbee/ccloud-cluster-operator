# Build the manager binary
FROM golang:1.17

ARG CCLOUD_EMAIL
ARG CCLOUD_PASSWORD

WORKDIR /workspace

# Copy the Go Modules manifests
COPY go.mod go.mod
COPY go.sum go.sum

# Copy the Confluent requisites
COPY confluent/cmd/confluent /bin/confluent
COPY confluent/script/.netrc /root

RUN go mod download

COPY main.go main.go
COPY api/ api/
COPY internal/ internal/
COPY controllers/ controllers/
COPY services/ services/

# Build
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -o runner main.go

WORKDIR /

USER root

RUN mkdir /manager

COPY confluent/script/setup /manager

RUN mv /workspace/runner /manager
RUN chmod 777 /manager/runner && chmod 777 /manager/setup


CMD ["/manager/runner","-D","FOREGROUND"]
