FROM golang:1.20.4-bullseye as build-backend

ARG SKIP_BACKEND_TEST=1

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

RUN go build -o app ./cmd/server

FROM node:19.9.0-alpine AS build-frontend

WORKDIR /build/frontend

COPY ./frontend/package*.json ./frontend/vue.config.js /build/frontend/
RUN npm install

COPY ./frontend/src /build/frontend/src

RUN npm run build

FROM golang:1.20.4-bullseye

WORKDIR /srv/webdict

EXPOSE ${PORT}

COPY --from=build-frontend /build/frontend/target/dist /srv/webdict/public
COPY --from=build-backend /build/backend/app /srv/webdict/app

RUN chown -R root:root /srv

CMD ["/srv/webdict/app"]