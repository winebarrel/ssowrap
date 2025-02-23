FROM golang:1.24 AS build
WORKDIR /src
COPY ./ /src/
ARG VERSION
RUN make BUILD_OPTS="-ldflags '-X main.version=${VERSION}'"

FROM debian:bookworm-slim

RUN <<EOF
set -e
apt-get update
apt-get install -y --no-install-recommends openssh-client
apt-get -y clean
rm -rf /var/lib/apt/lists/*
EOF

COPY --from=build /src/ssowrap /usr/local/bin/
