#
# Build Go binary inside base container.
#
FROM golang:1.16 as builder
WORKDIR /app
COPY . .
RUN make build


#
# Destination container.
#
FROM scratch
# Copy certificates and binary into the destination docker image.
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=builder /etc/passwd /etc/passwd
COPY --from=builder /app/bin/ethproxy /ethproxy
# Container settings.
ENV PORT 8080
EXPOSE 8080
USER nobody
ENTRYPOINT ["/ethproxy"]
