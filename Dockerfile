FROM --platform=$BUILDPLATFORM tonistiigi/xx AS xx

FROM --platform=$BUILDPLATFORM golang:1.20-alpine AS xbuild
WORKDIR /src
COPY --from=xx / /
ARG TARGETPLATFORM
ENV CGO_ENABLED=0
RUN --mount=target=. \
    --mount=type=cache,target=/root/.cache/go-build \
    --mount=type=cache,target=/go/pkg \
    xx-go build -trimpath -ldflags "-s -w" -o /dist/shopmon . && \
    xx-verify /dist/shopmon

FROM alpine:3.18
RUN apk add --no-cache ca-certificates
COPY --from=xbuild /dist/shopmon /usr/local/bin/shopmon
ENTRYPOINT ["/usr/local/bin/shopmon"]
