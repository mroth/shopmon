FROM --platform=$BUILDPLATFORM tonistiigi/xx AS xx

FROM --platform=$BUILDPLATFORM golang:1.18-alpine AS xbuild
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY internal internal
COPY --from=xx / /
ARG TARGETPLATFORM
ENV CGO_ENABLED=0
RUN xx-go build -trimpath -ldflags "-s -w" -o /dist/shopmon . && \
    xx-verify /dist/shopmon

FROM alpine:3.15
RUN apk add --no-cache ca-certificates
COPY --from=xbuild /dist/shopmon /usr/local/bin/shopmon
ENTRYPOINT ["/usr/local/bin/shopmon"]
