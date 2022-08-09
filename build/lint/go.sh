#!/bin/bash
set -Eeuo pipefail

util::infor "linting go files"

dirs=($(./pleasew query alltargets --include=go | grep -v third_party | cut -f1 -d":" | cut -c 3- | sort -u))
set +e
lint_output=$("${GO_LINT}" -set_exit_status ${dirs[@]} 2>&1)
lint_ec=$?
set -e
if [ $lint_ec -ne 0 ]; then
  printf "\n%s\n" "$lint_output"
  util::rerror "go files failed lint. To fix format errors, please run:
  $ ./pleasew run //build/fmt:go"
    exit 1
fi

util::rsuccess "linted go files"
