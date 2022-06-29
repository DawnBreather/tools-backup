FROM alpine

COPY --from=dutchcoders/transfer.sh /go/bin/transfersh /bin/transfersh
COPY --from=dutchcoders/transfer.sh /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt

ENTRYPOINT ["/bin/transfersh", "--listener", ":80"]
EXPOSE 80
