FROM golang:1.15 as builder

COPY . /src

WORKDIR /src
RUN CGO_ENABLED=0 go build -ldflags="-w -s"

FROM scratch

COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /src/helm3-oci-to-legacy-proxy /

EXPOSE 80
ENTRYPOINT [ "/helm3-oci-to-legacy-proxy" ]

LABEL org.opencontainers.image.source https://github.com/riksby/helm3-oci-to-legacy-proxy
