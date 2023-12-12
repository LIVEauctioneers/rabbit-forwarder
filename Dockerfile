FROM gcr.io/distroless/static

USER nonroot:nonroot

ADD --chown=nonroot:nonroot rabbit-forwarder /rabbit-forwarder

EXPOSE 80

CMD ["/rabbit-forwarder"]
