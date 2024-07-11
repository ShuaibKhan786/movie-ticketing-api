ARG GO_VERSION="1.22"
ARG ALPINE_VERSION="alpine3.20"


FROM golang:${GO_VERSION}-${ALPINE_VERSION} AS build-stage

WORKDIR /home

COPY go.mod go.sum .
RUN go mod download

COPY pkg/ ./pkg
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go


FROM scratch AS final-stage

WORKDIR /usr/local/bin


COPY --from=alpine:3.20 /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/


COPY --from=build-stage /home/server ./


CMD ["./server"]
