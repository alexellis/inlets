FROM golang:1.11-alpine as build

WORKDIR /go/src/github.com/alexellis/inlets

COPY .git               .git
COPY vendor             vendor
COPY pkg                pkg
COPY cmd                cmd
COPY main.go            .

ARG GIT_COMMIT
ARG VERSION
ARG OPTS

# add user in this stage because it cannot be done in next stage which is built from scratch
# in next stage we'll copy user and group information from this stage
RUN env ${OPTS} CGO_ENABLED=0 go build -ldflags "-s -w -X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION}" -a -installsuffix cgo -o /usr/bin/inlets \
    && addgroup -S app \
    && adduser -S -g app app

FROM scratch

COPY --from=build /etc/passwd /etc/group /etc/
COPY --from=build /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/
COPY --from=build /usr/bin/inlets /usr/bin/

USER app
EXPOSE 80

ENTRYPOINT ["/usr/bin/inlets"]
CMD ["--help"]
