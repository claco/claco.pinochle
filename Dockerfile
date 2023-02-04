###############################################################################
FROM --platform=${BUILDPLATFORM} golang:latest AS builder

ARG BUILDPLATFORM TARGETPLATFORM TARGETARCH TARGETOS
ARG DEBUG VCSREV=na VCSTAG=na VCSTIME=na VERSION=na

WORKDIR /usr/src/pinochle

ENV GOARCH=${TARGETARCH} GOOS=${TARGETOS}

RUN CGO_ENABLED=0 go install github.com/go-delve/delve/cmd/dlv@latest

COPY go.mod go.sum ./
RUN go mod download && go mod verify

COPY . .
RUN make build BINDIR=/usr/local/bin GOOS=${GOOS} GOARCH=${GOARCH} DEBUG=${DEBUG} VCSREV=${VCSREV} VCSTAG=${VCSTAG} VCSTIME=${VCSTIME} VERSION=${VERSION}


###############################################################################
FROM --platform=${BUILDPLATFORM} scratch AS debug

COPY --from=builder /go/bin/dlv /usr/local/bin/dlv
COPY --from=builder /usr/local/bin/pinochle /usr/local/bin/pinochle
COPY --from=builder /usr/src/pinochle/db ./db
COPY --from=builder /usr/src/pinochle /usr/src/pinochle
COPY ./files/etc /etc

USER pinochle:pinochle
ENV TMPDIR=/tmp

EXPOSE 2345

ENTRYPOINT ["dlv", "--listen=:2345", "--headless=true", "--api-version=2", "--accept-multiclient", "exec", "/usr/local/bin/pinochle", "--", "service", "run", "--debug", "--logs"]


###############################################################################
FROM --platform=${TARGETPLATFORM} scratch as pinochle

COPY --from=builder /usr/local/bin/pinochle /usr/local/bin/pinochle
COPY --from=builder /usr/src/pinochle/db ./db
COPY ./files/etc /etc

USER pinochle:pinochle
ENV TMPDIR=/tmp

HEALTHCHECK --interval=30s --timeout=2s --start-period=10s --retries=5 CMD [ "pinochle", "service", "check" ]

ENTRYPOINT [ "pinochle" ]
