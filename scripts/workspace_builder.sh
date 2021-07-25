#!/bin/bash
# This script prepares a Terraform Workspace with:
# * Terraform modules referenced by their absolute paths.
# * Terraform var files.
set -euo pipefail

mkdir -p "${OUTS}"

# modules configures modules for Terraform
# Terraform modules via Please work by determining the absolute path to
# the module source and updating the reference to that directory.
function modules {
    local abs_plz_out
    abs_plz_out="$(dirname "$PWD" | sed "s#$PKG##" | xargs dirname | xargs dirname)"

    for module in "${!MODULE_PATHS[@]}"; do
        path="${MODULE_PATHS[$module]}"
        find "${PKG_DIR}" -maxdepth 1 -name "*.tf" -exec sed -i "s#${module}#${abs_plz_out}/${path}#g" {} +
    done
}

# build_env_to_tf_srcs replaces various BUILD-time 
# environment variables in the Terraform source files.
# This is useful for re-using source file in multiple workspaces,
# such as templating a Terraform remote state configuration.
function build_env_to_tf_srcs {
    find "${OUTS}" -maxdepth 1 -name "*.tf" -exec sed -i "s#\$PKG#${PKG}#g" {} +
    find "${OUTS}" -maxdepth 1 -name "*.tf" -exec sed -i "s#\$PKG_DIR#${PKG_DIR}#g" {} +
    NAME="$(echo "${NAME}" | sed 's/^_\(.*\)_root$/\1/')"
    find "${OUTS}" -maxdepth 1 -name "*.tf" -exec sed -i "s#\$NAME#${NAME}#g" {} +
    find "${OUTS}" -maxdepth 1 -name "*.tf" -exec sed -i "s#\$ARCH#${ARCH}#g" {} +
    find "${OUTS}" -maxdepth 1 -name "*.tf" -exec sed -i "s#\$OS#${OS}#g" {} +
}

# auto_load_var_files copies the given var files into the 
# Terraform root and renames them so that they are auto-loaded 
# by Terraform so we don't have to use non-global `-var-file` flag.
function auto_load_var_files {
    SRCS_VAR_FILES=("${SRCS_VAR_FILES}")
    for i in "${!SRCS_VAR_FILES[@]}"; do
        var_file="${SRCS_VAR_FILES[i]}"
        cp "${var_file}" "${OUTS}/${i}-$(basename "${var_file}" | sed 's#\.tfvars#\.auto\.tfvars#')"
    done
}

# configure modules
if [[ -v SRCS_MODULES ]]; then
    modules
fi

# shift srcs into outs
for src in $SRCS_SRCS; do 
    cp "${src}" "${OUTS}/"
done

# substitute build env vars to srcs
build_env_to_tf_srcs

# shift var files into outs
if [[ -v SRCS_VAR_FILES ]]; then
    auto_load_var_files
fi
