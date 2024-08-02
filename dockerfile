ARG GO_VERSION="1.22"
ARG ALPINE_VERSION="alpine3.20"


FROM golang:${GO_VERSION}-${ALPINE_VERSION} AS build-stage

WORKDIR /home

COPY go.mod go.sum ./
COPY microservice/image-upload/go.mod ./microservice/image-upload/go.mod
RUN go mod download

COPY pkg/ ./pkg
COPY microservice/image-upload/pb ./microservice/image-upload/pb
COPY main.go ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server main.go


FROM scratch AS final-stage

WORKDIR /usr/local/bin

COPY --from=build-stage /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

COPY --from=build-stage /home/server ./


CMD ["./server"]
