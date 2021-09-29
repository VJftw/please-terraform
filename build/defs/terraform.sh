#!/bin/bash
#

set -Eeuo pipefail

# Bash version check
if [ -z "${BASH_VERSINFO}" ] || [ -z "${BASH_VERSINFO[0]}" ] || [ ${BASH_VERSINFO[0]} -lt 4 ]; then
    echo "This script requires Bash version >= 4"
    exit 1
fi

function log::debug {
    if [ -v PLZ_TF_DEBUG ]; then
        >&2 printf "debug> %s\n" "$1"
    fi
}

function _csv_to_array {
    local csv="$1"

    echo "${csv//,/ }"
}

function _parse_flag {
    local name="$1"
    shift
    while test $# -gt 0; do
        case "$1" in
            "--${name}="*)
                value=$(echo "$1" | cut -d= -f2-)
                echo "$value"
                log::debug "debug ${name}: $value (from: $@ )"
                shift
            ;;
            *)
                shift
            ;;
        esac
    done
}

# build_module
# This function prepares a Terraform module for use with Please by:
# * Replacing sub-modules (deps) with local references.
# * Ensuring all sub-modules have local references.
function build_module {
    local pkg name module_dir out url strip deps

    pkg="$(_parse_flag pkg "$@")"
    name="$(_parse_flag name "$@")"
    module_dir="$(_parse_flag module-dir "$@")"
    out="$(_parse_flag out "$@")"
    url="$(_parse_flag url "$@")"
    IFS=',' read -r -a strip <<< "$(_parse_flag strip "$@")"
    IFS=',' read -r -a deps <<< "$(_parse_flag deps "$@")"
    # strip=($(_parse_flag_array strip "$@"))
    # deps=($(_parse_flag_array deps "$@"))

    mv "${module_dir}" "${out}"

    # A Terraform module can have dependencies and can be depended on.
    # To accomodate this, we replace the deps aliases with the plz reference of dependency.
    mkdir "${out}/modules/"
    for m in "${deps[@]}"; do
        replace=$(basename "$m")
        mapfile -t searches <"${m}/.module_aliases"
        for search in "${searches[@]}"; do
            find . -name "*.tf" -exec sed -i  "s#\"[^\"]*${search}[^\"]*\"#\"./modules/${replace}\"#g" {} +
        done
        cp -r "$m" "${out}/modules/"
    done

    # strip directories
    for s in "${strip[@]}"; do
        rm -rf "${out:?}/${s}"
    done

    # Add aliases to the module for dependants to use.
    # add a replace-me search for an interesting part of the URL
    echo "${url}" | cut -f3-5 -d/ > "${out}/.module_aliases"
    # add a replace-me search for the canonical Please build rule
    echo "${pkg}:${name}" >> "${out}/.module_aliases"
}

# build_workspace
# This function prepares a Terraform Workspace with:
# * Terraform modules referenced by their absolute paths.
# * Terraform var files.
function build_workspace {
    local pkg name os arch terraform_binary out pkg_dir srcs var_files module_paths

    pkg="$(_parse_flag pkg "$@")"
    name="$(_parse_flag name "$@")"
    os="$(_parse_flag os "$@")"
    arch="$(_parse_flag arch "$@")"
    out="$(_parse_flag out "$@")"
    pkg_dir="$(_parse_flag pkg-dir "$@")"
    IFS=',' read -r -a srcs <<< "$(_parse_flag srcs "$@")"
    IFS=',' read -r -a var_files <<< "$(_parse_flag var_files "$@")"
    IFS=',' read -r -a module_paths <<< "$(_parse_flag module-paths "$@")"

    mkdir -p "${out}"

    # configure modules
    # Terraform modules via Please work by determining the absolute path to
    # the module source and updating the reference to that directory.
    local abs_plz_out
    abs_plz_out="$(dirname "$PWD" | sed "s#$pkg##" | xargs dirname | xargs dirname)"
    for module_path in "${module_paths[@]}"; do
        # module="$(echo "$module_path" | cut -f1 -d=)"
        # path="$(echo "$module_path" | cut -f2 -d=)"
        module="${module_path%=*}"
        path="${module_path#*=}"
        find "${pkg_dir}" -maxdepth 1 -name "*.tf" -exec sed -i "s#${module}#${abs_plz_out}/${path}#g" {} +
    done

    # shift srcs into outs
    for src in "${srcs[@]}"; do
        cp "${src}" "${out}/"
    done

    # substitute build env vars to srcs
    # This is useful for re-using source file in multiple workspaces,
    # such as templating a Terraform remote state configuration.
    find "${out}" -maxdepth 1 -name "*.tf" -exec sed -i "s#\$PKG#${pkg}#g" {} +
    find "${out}" -maxdepth 1 -name "*.tf" -exec sed -i "s#\$PKG_DIR#${pkg_dir}#g" {} +
    NAME="$(echo "${name}" | sed 's/^_\(.*\)_root$/\1/')"
    find "${out}" -maxdepth 1 -name "*.tf" -exec sed -i "s#\$NAME#${name}#g" {} +
    find "${out}" -maxdepth 1 -name "*.tf" -exec sed -i "s#\$ARCH#${arch}#g" {} +
    find "${out}" -maxdepth 1 -name "*.tf" -exec sed -i "s#\$OS#${os}#g" {} +

    # shift var files into outs
    # copies the given var files into the
    # Terraform root and renames them so that they are auto-loaded
    # by Terraform so we don't have to use non-global `-var-file` flag.
    for i in "${!var_files[@]}"; do
        var_file="${var_files[i]}"
        cp "${var_file}" "${out}/${i}-$(basename "${var_file}" | sed 's#\.tfvars#\.auto\.tfvars#')"
    done
}

# run
# This function runs Terraform in the target's working directory with the following features:
# - Copying plugins into a plugin cache directory.
# - Strip out various noisy output (https://github.com/hashicorp/terraform/issues/20960)
# - Executing pre-binaries.
# - Executing post-binaries.
# - Executing terraform with user provided flags.
function run {

    terraform_binary="$(_parse_flag terraform-binary "$@")"
    os="$(_parse_flag os "$@")"
    arch="$(_parse_flag arch "$@")"
    terraform_root="$(_parse_flag terraform-root "$@")"
    IFS=',' read -r -a plugins <<< "$(_parse_flag plugins "$@")"
    IFS=',' read -r -a pre_binaries <<< "$(_parse_flag pre-binaries "$@")"
    IFS=',' read -r -a post_binaries <<< "$(_parse_flag post-binaries "$@")"
    IFS=' ' read -r -a extra_args <<< "$(_parse_flag extra-args "$@")"
    # extra_args="$(_parse_flag extra-args "$@")"

    # determine absolute path to repository root.
    repo_root="${PWD}"

    # use absolute paths for the terraform root and terraform binary.
    terraform_root="${repo_root}/${terraform_root}"
    terraform_binary="${repo_root}/${terraform_binary}"

    # determine the minor version of Terraform.
    terraform_minor_version="$(head -n1 < <($terraform_binary version) | awk '{ print $2 }' | cut -f1-2 -d\.)"
    # configure where to place Terraform plugins
    plugin_dir="${terraform_root}/_plugins"
    mkdir -p "${plugin_dir}"

    # set some Terraform options.
    export TF_PLUGIN_CACHE_DIR="${plugin_dir}"

    # enable the use of `terraform` in pre/post binaries.
    PATH="$(dirname "${terraform_binary}"):$PATH"
    export PATH

    # We cannot run Terraform commands in the `plz-out/gen/<rule>` workspace
    # as Terraform creates symlinks which plz warns us may be removed, thus
    # we create a `plz-out/terraform` directory and `rsync` the following:
    # - Generated Terraform Root files.
    # - Terraform plugins into a rule-local cache directory (`plz-out/terraform/<rule>/_plugins`).
    # Terraform modules are referenced absolutely to their `plz-out/gen/<module rule>` counterparts.
    terraform_workspace="${terraform_root//plz-out\/gen/plz-out\/terraform}"
    mkdir -p "${terraform_workspace}"
    rsync -ah --delete --exclude=.terraform* --exclude=*.tfstate* --exclude=/_plugins "${terraform_root}/" "${terraform_workspace}/"

    # copy plugins (providers)
    if ((${#plugins[@]})); then
        case "${terraform_minor_version}" in
            "v0.11") run_plugins_v0.11+ "$plugin_dir" "$os" "$arch" "${plugins[@]}";;
            "v0.12") run_plugins_v0.11+ "$plugin_dir" "$os" "$arch" "${plugins[@]}";;
            "v0.13") run_plugins_v0.13+ "$plugin_dir" "$os" "$arch" "${plugins[@]}";;
            *) run_plugins_v0.13+ "$plugin_dir" "$os" "$arch" "${plugins[@]}";;
        esac
    fi
    cd "${terraform_workspace}"

    # execute pre_binaries
    for bin in "${pre_binaries[@]}"; do
        "${repo_root}/${bin}"
    done

    # execute terraform with the given cmd
    run_tf_clean_output "${terraform_binary}" "${extra_args[@]}"

    # execute post_binaries
    for bin in "${post_binaries[@]}"; do
        "${repo_root}/${bin}"
    done

}

# run_plugins_v0.11+ configures plugins for Terraform 0.11+
# Terraform v0.11+ store plugins in the following structure:
# `./${os}_${arch}/${binary}`
# e.g. ``./linux_amd64/terraform-provider-null_v2.1.2_x4`
function run_plugins_v0.11+ {
    local plugin_base_dir="$1"; shift
    local os="$1"; shift
    local arch="$1"; shift
    local plugin_bin
    local plugins=("$@")
    plugin_dir="${plugin_base_dir}/${os}_${arch}"
    mkdir -p "${plugin_dir}"
    for plugin in "${plugins[@]}"; do
        plugin_bin="$(find "$plugin" -not -path '*/\.*' -type f | head -n1)"
        rsync -ah "$plugin_bin" "${plugin_dir}/"
    done
}


# run_plugins_v0.13+ configures plugins for Terraform 0.13+
# Terraform v0.13+ store plugins in the following structure:
# `./${registry}/${namespace}/${type}/${version}/${os}_${arch}/${binary}`
# e.g. `./registry.terraform.io/hashicorp/null/2.1.2/linux_amd64/terraform-provider-null_v2.1.2_x4`
function run_plugins_v0.13+ {
    local plugins registry namespace provider_name version plugin_dir plugin_bin
    local plugin_base_dir="$1"; shift
    local os="$1"; shift
    local arch="$1"; shift
    plugins=("$@")
    for plugin in "${plugins[@]}"; do
        registry=$(<"${plugin}/.registry")
        namespace=$(<"${plugin}/.namespace")
        provider_name=$(<"${plugin}/.provider_name")
        version=$(<"${plugin}/.version")
        plugin_dir="${plugin_base_dir}/${registry}/${namespace}/${provider_name}/${version}/${os}_${arch}"
        plugin_bin="$(find "$plugin" -not -path '*/\.*' -type f | head -n1)"

        mkdir -p "${plugin_dir}"
        rsync -ah "$plugin_bin" "${plugin_dir}/"
    done
}

# run_tf_clean_output strips the Terraform output down.
# This is useful in CI/CD where Terraform logs are usually noisy by default.
function run_tf_clean_output {
    local terraform_binary args
    terraform_binary="$1"; shift
    args=("${@}")

    # set TF_CLEAN_OUTPUT to false by default.
    export TF_CLEAN_OUTPUT="${TF_CLEAN_OUTPUT:-false}"

    echo "..> terraform ${args[*]}"
    if [ "${TF_CLEAN_OUTPUT}" == "false" ]; then
        "${terraform_binary}" "${args[@]}"
    else
        "${terraform_binary}" "${args[@]}" \
        | sed '/successfully initialized/,$d' \
        | sed "/You didn't specify an \"-out\"/,\$d" \
        | sed '/.terraform.lock.hcl/,$d' \
        | sed '/Refreshing state/d' \
        | sed '/The refreshed state will be used to calculate this plan/d' \
        | sed '/persisted to local or remote state storage/d' \
        | sed '/^[[:space:]]*$/d'
    fi
}

case "$1" in
    "build_module")
        build_module "$@"
    ;;
    "build_workspace")
        build_workspace "$@"
    ;;
    "run")
        run "$@"
    ;;
esac
