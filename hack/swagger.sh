#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

pushd "${base_dir}" >/dev/null

if ! [ -e "${base_dir}/bin/swag" ]; then
    export GOBIN="${base_dir}/bin"
    echo "Installing swag tool"
    go install github.com/swaggo/swag/cmd/swag@latest

    echo ""
fi

export swag="${base_dir}/bin/swag"

echo "Generating Swagger documentation"
"${swag}" init -g api.go --dir pkg/server/api/v1/,pkg/server/storage/ -o ./ --ot yaml

echo ""

echo "Formatting swagger comments"
"${swag}" fmt -g api.go --dir pkg/server/api/v1/

popd >/dev/null
