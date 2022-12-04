FROM golang:1.18-bullseye as build-backend

ARG SKIP_BACKEND_TEST

ADD backend /build/backend
WORKDIR /build/backend

RUN echo go version: `go version`

ENV GOFLAGS="-mod=vendor"

# run tests
RUN \
  if [ -z "$SKIP_BACKEND_TEST" ] ; then \
    go test -race -p 1 -timeout="${BACKEND_TEST_TIMEOUT:-300s}" -covermode=atomic ./... ; \
  else \
    echo "skip backend tests and linter" \
  ; fi

RUN go build -o webdict ./cmd/server

FROM golang:1.18-bullseye

WORKDIR /srv

COPY --from=build-backend /build/backend/webdict /srv/webdict

CMD ["/srv/webdict"]