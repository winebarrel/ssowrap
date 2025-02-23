FROM golang:1.24 AS build

WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download
COPY Makefile *.go /src/
COPY cmd /src/cmd
ARG VERSION
RUN make BUILD_OPTS="-ldflags '-X main.version=$VERSION'"

FROM debian:bookworm-slim

RUN <<EOF
set -e
apt-get update
apt-get install -y --no-install-recommends \
  ca-certificates \
  curl \
  openssh-client \
  awscli
apt-get -y clean
rm -rf /var/lib/apt/lists/*

# https://docs.aws.amazon.com/systems-manager/latest/userguide/install-plugin-debian-and-ubuntu.html
case $(uname -m) in
x86_64)
  SESSION_MANAGER_PLUGIN_ARCH=ubuntu_64bit
  ;;
aarch64)
  SESSION_MANAGER_PLUGIN_ARCH=ubuntu_arm64
  ;;
esac

if [ -n "$SESSION_MANAGER_PLUGIN_ARCH" ]; then
  curl -O https://s3.amazonaws.com/session-manager-downloads/plugin/latest/$SESSION_MANAGER_PLUGIN_ARCH/session-manager-plugin.deb
  dpkg -i session-manager-plugin.deb
  rm session-manager-plugin.deb
fi
EOF

COPY --from=build /src/ssowrap /usr/local/bin/
