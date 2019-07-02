FROM golang:alpine as builder
ADD . /build/
WORKDIR /build
RUN \
    apk add --no-cache \
        ca-certificates \
        gcc \
        git \
        musl-dev \
    && go get . \
    && go test -v . \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -ldflags '-extldflags "-static"' -o ecr_exporter .

FROM scratch
COPY --from=builder /build/ecr_exporter /
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 9606/tcp
USER 65534
ENTRYPOINT ["/ecr_exporter"]
