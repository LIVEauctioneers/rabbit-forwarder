FROM liveauctioneers/alpine

RUN apk add --update ca-certificates tzdata && \
    rm -rf /var/cache/apk/* /tmp/*

WORKDIR /app/

COPY ./rabbit-amazon-forwarder /app/rabbit-amazon-forwarder

EXPOSE 80

CMD ["./rabbit-amazon-forwarder"]