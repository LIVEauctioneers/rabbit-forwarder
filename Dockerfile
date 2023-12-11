FROM gcr.io/distroless/static

USER nonroot:nonroot

ADD --chown=nonroot:nonroot rabbit-amazon-forwarder /rabbit-amazon-forwarder

EXPOSE 80

CMD ["/rabbit-amazon-forwarder"]
