FROM --platform=$BUILDPLATFORM tonistiigi/xx AS xx

FROM --platform=$BUILDPLATFORM golang:1.18-alpine AS xbuild
COPY --from=xx / /
ARG TARGETPLATFORM
ENV CGO_ENABLED=0
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY internal internal
RUN xx-go build -o /dist/shopmon . && \
    xx-verify /dist/shopmon

FROM alpine:3.15
RUN apk add --no-cache ca-certificates
COPY --from=xbuild /dist/shopmon /usr/local/bin/shopmon
ENTRYPOINT ["/usr/local/bin/shopmon"]
