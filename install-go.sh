#!/bin/sh

WORKDIR=$(pwd)
CACHE_PATH="${WORKDIR}/.cache"

GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)
GO_VERSION=$(cat "${WORKDIR}/go-version")
GO_TAR="go${GO_VERSION}.${GOOS}-${GOARCH}.tar.gz"
GO_SRC="https://go.dev/dl/${GO_TAR}"
GO_CUSTOM_PATH="${GO_CUSTOM_PATH:-${WORKDIR}/go}"
DL_PATH="${CACHE_PATH}/downloads/${GO_TAR}"

echo "GOOS: ${GOOS}"
echo "GOARCH: ${GOARCH}"
echo "GO_CUSTOM_PATH: ${GO_CUSTOM_PATH}"

usage() {
    echo
    echo "To use the installed go binary, add these lines to your .bashrc or .zshrc file."
    echo "      export PATH=\"${GO_CUSTOM_PATH}/go/bin:\$PATH\""
    echo "      export PATH=\"${GO_CUSTOM_PATH}/bin:\$PATH\""
    echo "      export GOROOT=\"${GO_CUSTOM_PATH}/go\""
    echo "      export GOPATH=\"$(dirname $GO_CUSTOM_PATH)/go\""
}

download_go(){
    if [ ! -e "$DL_PATH" ]; then
        mkdir -p "$(dirname $DL_PATH)" && \
            wget --progress=bar:force:noscroll -O "${DL_PATH}" "${GO_SRC}"
    fi
}

if [ -f "${GO_CUSTOM_PATH/go-version}" ] && [ "$GO_VERSION" = "$(cat $GO_CUSTOM_PATH/go-version)" ]; then
    echo "Go is already installed" && usage
    exit 0
else
    echo "Downloading ${GO_SRC}..." && \
        download_go && \
        echo "Extracting ${GO_TAR} to ${GO_CUSTOM_PATH}..." && \
        rm -rf ${GO_CUSTOM_PATH} && mkdir -p ${GO_CUSTOM_PATH} && \
        tar -C ${GO_CUSTOM_PATH} -xzf "${DL_PATH}" && \
        echo $GO_VERSION > "$GO_CUSTOM_PATH/go-version" && \
        echo "Installed Go ${GO_VERSION} to ${GO_CUSTOM_PATH}" && usage
fi

cd $WORKDIR
