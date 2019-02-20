FROM golang:1.10 as build

WORKDIR /go/src/github.com/alexellis/inlets

COPY .git             .git
COPY vendor             vendor
COPY pkg                pkg
COPY main.go            .
COPY parse_upstream.go  .

ARG GIT_COMMIT
ARG VERSION

RUN CGO_ENABLED=0 go build -ldflags "-X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION}" -a -installsuffix cgo -o /usr/bin/inlets

FROM alpine:3.9
RUN apk add --force-refresh ca-certificates

COPY --from=build /usr/bin/inlets /root/

EXPOSE 80

WORKDIR /root/

CMD ["/usr/bin/inlets"]
