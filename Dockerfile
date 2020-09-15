FROM golang:alpine as golang
WORKDIR /app
COPY . .
RUN apk add --no-cache git libc-dev musl-dev build-base gcc ca-certificates \
    && export VCS_REF=$(git tag --points-at HEAD) \
    && echo $VCS_REF \
    && go get -v . \
    && go install -v github.com/rakyll/statik \
    && $GOPATH/bin/statik -src="./www/" \
    && CGO_ENABLED=1 go build -ldflags "-s -w -X main.Version=$VCS_REF" .

FROM alpine
COPY --from=golang /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=golang /app/mantle /app/mantle

VOLUME /data
ENTRYPOINT ["/app/mantle", "--port", "8000", "--config", "/data/config.json"]
