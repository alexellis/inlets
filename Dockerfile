FROM golang:1.15-alpine as build

ENV GO111MODULE=on
ENV CGO_ENABLED=0

WORKDIR /go/src/github.com/inlets/inlets

COPY vendor             vendor
COPY go.mod             .
COPY go.sum             .
COPY pkg                pkg
COPY cmd                cmd
COPY main.go            .

ARG GIT_COMMIT
ARG VERSION
ARG OPTS

ARG TARGETARCH

RUN test -z "$(gofmt -l $(find . -type f -name '*.go' -not -path "./vendor/*"))" || { echo "Run \"gofmt -s -w\" on your Golang code"; exit 1; }
RUN go test -mod=vendor $(go list ./...) -cover

# add user in this stage because it cannot be done in next stage which is built from scratch
# in next stage we'll copy user and group information from this stage
RUN env ${OPTS} go build -mod=vendor \
    -ldflags "-s -w -X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION}" \
    -a -installsuffix cgo -o /usr/bin/inlets

RUN addgroup -S app \
    && adduser -S -g app app

FROM scratch

COPY --from=build /etc/passwd /etc/group /etc/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/bin/inlets /usr/bin/

EXPOSE 80

VOLUME /tmp/

ENTRYPOINT ["/usr/bin/inlets"]
CMD ["--help"]
