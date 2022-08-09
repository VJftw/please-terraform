#!/bin/bash
set -Eeuo pipefail

util::infor "formatting go files"

dirs=($(./pleasew query alltargets --include=go | grep -v third_party | cut -f1 -d":" | cut -c 3- | sort -u))
"${GO_FMT}" -s -w ${dirs[@]}

"${GO_BIN}" mod tidy

util::rsuccess "formatted go files"
