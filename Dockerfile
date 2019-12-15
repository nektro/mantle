FROM golang:alpine as golang
WORKDIR /go/src/mantle
COPY . .
RUN CGO_ENABLED=0 go install -ldflags '-extldflags "-static"'

FROM scratch
COPY --from=golang /go/bin/mantle /app
COPY --from=golang /go/src/mantle/public /public

ENTRYPOINT ["/app"]
