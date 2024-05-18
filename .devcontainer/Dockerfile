ARG VARIANT="0-1.22-bookworm"

FROM mcr.microsoft.com/vscode/devcontainers/go:${VARIANT}

RUN apt-get update && export DEBIAN_FRONTEND=noninteractive && apt-get install -y protobuf-compiler

USER vscode

# Install Go Protobuf tools.
RUN go install github.com/go-delve/delve/cmd/dlv@latest
RUN go install github.com/spf13/cobra-cli@latest