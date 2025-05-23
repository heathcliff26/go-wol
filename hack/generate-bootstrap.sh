#!/bin/bash

set -e

base_dir="$(dirname "${BASH_SOURCE[0]}" | xargs realpath)/.."

bootstrap_file="${base_dir}/static/bootstrap/bootstrap.css"
index_html="${base_dir}/static/index.html"
utils_js="${base_dir}/static/js/utils.js"
output_file="${base_dir}/static/css/bootstrap.css"

if ! command -v purgecss >/dev/null 2>&1; then
    echo "purgecss is not installed, trying to install it"
    npm install -g purgecss
fi

echo "Creating trimmed boostrap css"
purgecss --css "${bootstrap_file}" --content "${index_html}" --content "${utils_js}" -o "${output_file}" \
    -s "alert-success" -s "alert-warning" -s "alert-danger"
