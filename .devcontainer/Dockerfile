# See https://github.com/devcontainers/images/tree/main/src/go

# [Choice] Go version: 1, 1.27, 1-bookworm, 1.27-bookworm, ...
ARG VARIANT="1-bookworm"
FROM mcr.microsoft.com/devcontainers/go:${VARIANT}

# [Optional] Uncomment this section to install additional OS packages.
# hadolint ignore=DL3008,DL3009
RUN apt-get update && export DEBIAN_FRONTEND=noninteractive && apt-get -y upgrade \
    && apt-get -y install --no-install-recommends libcurl4-openssl-dev

# [Optional] Uncomment the next lines to use go get to install anything else you need
USER vscode
RUN go install mvdan.cc/gofumpt@latest

# [Optional] Uncomment this line to install global node packages.
# RUN su vscode -c "source /usr/local/share/nvm/nvm.sh && npm install -g <your-package-here>" 2>&1
