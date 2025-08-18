FROM golang:1.25.0-bookworm AS build

WORKDIR /src
COPY go.mod go.sum /src/
RUN go mod download
COPY Makefile *.go /src/
COPY cmd /src/cmd
ARG VERSION
RUN make

FROM debian:bookworm-slim

RUN --mount=type=cache,target=/var/apt/cache <<EOF
set -e
rm /etc/apt/apt.conf.d/docker-clean
apt-get update -qq
apt-get install -qq -y --no-install-recommends \
  ca-certificates \
  curl \
  openssh-client \
  awscli

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
  curl -sSfO https://s3.amazonaws.com/session-manager-downloads/plugin/latest/$SESSION_MANAGER_PLUGIN_ARCH/session-manager-plugin.deb
  dpkg -i session-manager-plugin.deb
  rm session-manager-plugin.deb
fi
EOF

COPY --from=build /src/ssowrap /usr/local/bin/
