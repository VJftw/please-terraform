#!/bin/bash
# This script runs Terraform in the target's working directory with the following features:
# - Copying plugins into a plugin cache directory.
# - Strip out various noisy output (https://github.com/hashicorp/terraform/issues/20960)
# - Executing pre-binaries.
# - Executing post-binaries.
# - Executing terraform with user provided flags.
set -euo pipefail

# ABS is the absolute path to the repository.
ABS="${PWD}"
# PLZ_OUT is the absolute path to the plz-out directory.
PLZ_OUT="${ABS}/${PLZ_OUT_RELATIVE}"
# TERRAFORM_WORKSPACE is the absolute path to the runtime Terraform workspace that we run Terraform commands in.
TERRAFORM_WORKSPACE="${PLZ_OUT}/terraform/${TERRAFORM_ROOT//$PLZ_OUT\/gen/}"
# TERRAFORM_MINOR_VERSION is the version of Terraform to the minor version. e.g. 0.11, 0.12, 0.13, 0.14, ...
TERRAFORM_MINOR_VERSION="$(head -n1 < <($TERRAFORM_BIN version) | awk '{ print $2 }' | cut -f1-2 -d\.)"
# PLUGIN_DIR is the absolute path to the Terraform plugin cache to prepare.
PLUGIN_DIR="${TERRAFORM_WORKSPACE}/_plugins"
# TERRAFORM_BIN is the absolute path to the Terraform binary.
TERRAFORM_BIN="${ABS}/${TERRAFORM_BIN}"
# TF_CLEAN_OUTPUT is whether or not to clean the Terraform CLI output on a best effort basis. This is incompatible with interactive workflows.
TF_CLEAN_OUTPUT="${TF_CLEAN_OUTPUT:-false}"

export TF_PLUGIN_CACHE_DIR="${PLUGIN_DIR}"
PATH="$(dirname "${TERRAFORM_BIN}"):$PATH"
export PATH

# prepare_workspace prepares a directory to run Terraform commands in.
# We cannot run Terraform commands in the `plz-out/gen/<rule>` workspace
# as Terraform creates symlinks which plz warns us may be removed, thus
# we create a `plz-out/terraform` directory and `rsync` the following:
# - Generated Terraform Root files.
# - Terraform plugins into a rule-local cache directory (`plz-out/terraform/<rule>/_plugins`).
# Terraform modules are referenced absolutely to their `plz-out/gen/<module rule>` counterparts.
function prepare_workspace {
    local terraform_workspace
    terraform_workspace="$1"
    mkdir -p "${terraform_workspace}"
    rsync -ah --delete --exclude=/.terraform --exclude=/_plugins "${TERRAFORM_ROOT}/" "${terraform_workspace}/"

    # copy plugins (providers)
    if [[ -v PLUGINS ]]; then
        case "${TERRAFORM_MINOR_VERSION}" in
            "v0.11") plugins_v0.11+ ;;
            "v0.12") plugins_v0.11+ ;;
            "v0.13") plugins_v0.13+ ;;
            *) plugins_v0.13+ ;;
        esac
    fi
}

# plugins_v0.11+ configures plugins for Terraform 0.11+
# Terraform v0.11+ store plugins in the following structure:
# `./${os}_${arch}/${binary}`
# e.g. ``./linux_amd64/terraform-provider-null_v2.1.2_x4`
function plugins_v0.11+ {
    local plugin_dir
    local plugin_bin
    plugin_dir="${PLUGIN_DIR}/${OS}_${ARCH}"
    mkdir -p "${plugin_dir}"
    for plugin in "${PLUGINS[@]}"; do
        plugin_bin="$(find "$plugin" -not -path '*/\.*' -type f | head -n1)"
        rsync -ah "$plugin_bin" "${plugin_dir}/"
    done
}

# plugins_v0.13+ configures plugins for Terraform 0.13+
# Terraform v0.13+ store plugins in the following structure:
# `./${registry}/${namespace}/${type}/${version}/${os}_${arch}/${binary}`
# e.g. `./registry.terraform.io/hashicorp/null/2.1.2/linux_amd64/terraform-provider-null_v2.1.2_x4`
function plugins_v0.13+ {
    local registry namespace provider_name version plugin_dir plugin_bin
    for plugin in "${PLUGINS[@]}"; do
        registry=$(<"${plugin}/.registry")
        namespace=$(<"${plugin}/.namespace")
        provider_name=$(<"${plugin}/.provider_name")
        version=$(<"${plugin}/.version")
        plugin_dir="${PLUGIN_DIR}/${registry}/${namespace}/${provider_name}/${version}/${OS}_${ARCH}"
        plugin_bin="$(find "$plugin" -not -path '*/\.*' -type f | head -n1)"
        mkdir -p "${plugin_dir}"
        rsync -ah "$plugin_bin" "${plugin_dir}/"
    done
}

# tf_clean_output strips the Terraform output down.
# This is useful in CI/CD where Terraform logs are usually noisy by default.
function tf_clean_output {
    local cmds extra_args is_last
    IFS=" " read -r -a cmds <<< "$1"
    shift
    is_last="$1"
    shift
    extra_args=("${@}")

    args=("${cmds[@]}")
    if [ "${is_last}" == "true" ]; then
        args=("${args[@]}" "${extra_args[@]}")
    fi
    echo "..> terraform ${args[*]}"
    if [ "${TF_CLEAN_OUTPUT}" == "false" ]; then
        "${TERRAFORM_BIN}" "${args[@]}"
    else
        "${TERRAFORM_BIN}" "${args[@]}" \
        | sed '/successfully initialized/,$d' \
        | sed "/You didn't specify an \"-out\"/,\$d" \
        | sed '/.terraform.lock.hcl/,$d' \
        | sed '/Refreshing state/d' \
        | sed '/The refreshed state will be used to calculate this plan/d' \
        | sed '/persisted to local or remote state storage/d' \
        | sed '/^[[:space:]]*$/d'
    fi
}

prepare_workspace "${TERRAFORM_WORKSPACE}"
cd "${TERRAFORM_WORKSPACE}"

# execute pre_binaries
for bin in "${PRE_BINARIES[@]}"; do
    "${ABS}/${bin}"
done

# execute terraform cmds
if [[ -v TERRAFORM_CMDS ]]; then
    for i in "${!TERRAFORM_CMDS[@]}"; do
        cmd="${TERRAFORM_CMDS[i]}"
        if [ $((i+1)) == "${#TERRAFORM_CMDS[@]}" ]; then
            tf_clean_output "${cmd}" "true" "$@"
        else
            tf_clean_output "${cmd}" "false" "$@"
        fi

        echo ""
    done
else
    # if there's no TERRAFORM_CMDS given, we assume that we just want to run Terraform directly with the given args.
    echo "..> terraform ${@}"
    "${TERRAFORM_BIN}" "${@}"
fi

# execute post_binaries
for bin in "${POST_BINARIES[@]}"; do
    "${ABS}/${bin}"
done
