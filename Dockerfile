FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS builder
ARG APP_NAME=111111111
RUN echo "Build app [${APP_NAME}]"

WORKDIR /build
ENV GOOS linux
ENV CGO_ENABLED 0
ENV GOPATH /go
ENV GOCACHE /go-build
COPY go.mod go.sum ./
RUN --mount=type=cache,target=/go/pkg/mod/cache go mod download
#COPY cmd/${APP_NAME} pkg/ internal/ ./
COPY cmd cmd
#COPY pkg pkg
COPY internal internal
RUN --mount=type=cache,target=/go/pkg/mod/cache \
    --mount=type=cache,target=/go-build \
    go build -o app ./cmd/${APP_NAME}

FROM alpine:3.17.3 as final
RUN adduser -D -u 1000 appuser

RUN apk --no-cache add tzdata
WORKDIR /dist
COPY --from=builder /build/app .
USER appuser
CMD ["/dist/app"]
