FROM golang:1.10 as build

WORKDIR /go/src/github.com/alexellis/inlets

COPY .git               .git
COPY vendor             vendor
COPY pkg                pkg
COPY main.go            .
COPY Makefile           .
COPY parse_upstream.go  .

ARG GIT_COMMIT
ARG VERSION

RUN make install

FROM alpine:3.9
RUN apk add --force-refresh ca-certificates

# Add non-root user
RUN addgroup -S app && adduser -S -g app app \
  && mkdir -p /home/app || : \
  && chown -R app /home/app

COPY --from=build /usr/local/bin/inlets /usr/bin/
WORKDIR /home/app

USER app
EXPOSE 80

ENTRYPOINT ["inlets"]
CMD ["-help"]
