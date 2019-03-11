FROM golang:1.10 as build

WORKDIR /go/src/github.com/alexellis/inlets

COPY .git               .git
COPY vendor             vendor
COPY pkg                pkg
COPY cmd                cmd
COPY main.go            .

ARG GIT_COMMIT
ARG VERSION

RUN CGO_ENABLED=0 go build -ldflags "-s -w -X main.GitCommit=${GIT_COMMIT} -X main.Version=${VERSION}" -a -installsuffix cgo -o /usr/bin/inlets

FROM alpine:3.9
RUN apk add --force-refresh ca-certificates

# Add non-root user
RUN addgroup -S app && adduser -S -g app app \
  && mkdir -p /home/app || : \
  && chown -R app /home/app

COPY --from=build /usr/bin/inlets /usr/bin/
WORKDIR /home/app

USER app
EXPOSE 80

ENTRYPOINT ["inlets"]
CMD ["-help"]
