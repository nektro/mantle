FROM alpine
WORKDIR /app
COPY ./bin/mantle-linux-amd64 /app/mantle
VOLUME /data
ENTRYPOINT ["/app/mantle", "--port", "8000", "--config", "/data/config.json", "--skip-translation-fetch"]
