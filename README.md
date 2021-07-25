# please-terraform
Terraform integration w/ the Please build system.

This includes support for the following:
 * `terraform_provider`: Terraform Providers
 * `terraform_module`: Terraform Remote Modules
 * `terraform_module`: Terraform Local Modules
 * `terraform_toolchain`: Multiple versions of Terraform
 * Terraform fmt/validate


## `terraform_toolchain`

This build rule allows you to specify a Terraform version to download and re-use in `terraform_root` rules. You can repeat this for multiple versions if you like, see `//third_party/terraform/BUILD` for examples.

## `terraform_provider`

This build rule allows you to specify a [Terraform provider](https://www.terraform.io/docs/providers/index.html) to re-use in your `terraform_root` rules. See `//third_party/terraform/provider/BUILD` for examples.

## `terraform_module`

This build rule allows you to specify a [Terraform module](https://www.terraform.io/docs/language/modules/index.html) to re-use in your `terraform_root` rules or as dependencies in other `terraform_module` rules. Terraform modules can be sourced remotely or exist locally on the filesystem. 

See `//third_party/terraform/module/BUILD` for examples of remote Terraform modules.
See `//terraform/examples/<version>/my_module/BUILD` for examples of local terraform modules.

In your Terraform source code, you should refer to your modules by their canonical build label. e.g.:

```
module "remote_module" {
    source = "//third_party/terraform/module:a_module"
}

module "my_module" {
    source = "//terraform/examples/0.12/my_module:my_module"
}
``` 

If your module has providers or required providers configuration, you must include them as deps.

## `terraform_root`

This build rule allows to specify a [Terraform root module](https://www.terraform.io/docs/language/modules/index.html#the-root-module) which is the root configuration where Terraform will be executed. In this build rule, you reference the `srcs` for the root module as well as the providers and modules those `srcs` use. 

Terraform Providers are mirrored into a local directory for Terraform to source them from (https://www.terraform.io/docs/cli/config/config-file.html#explicit-installation-method-configuration).


We support substitution of the following please build environment variables into your source terraform files:
 - `PKG`
 - `PKG_DIR`
 - `NAME`
 - `ARCH`
 - `OS` 
This allows you to template Terraform code to keep your code DRY. for example: A terraform remote state configuration can that can be re-used in all `terraform_root`s:
```
terraform {
  backend "s3" {
    region         = "eu-west-1"
    bucket         = "my-terraform-state"
    key            = "$PKG/$NAME.tfstate"
    dynamodb_table = "my-terraform-state-lock"
    encrypt        = true
  }
}
```
The above will result in a terraform state tree consistent with the structure of your repository.

This build rule generates the following subrules which perform the Terraform workflows:
 * `_plan`
 * `_apply`
 * `_destroy`
 * `_bin` for all other workflows e.g. `plz run //my_infrastructure_tf_bin -- init && plz run //my_infrastructure_tf_bin -- console`

For all of these workflows, we support passing in flags via please as expected, e.g.:
```
$ plz run //my_tf:my_tf_plan -- -lock=false
$ plz run //my_tf:my_tf_import -- resource_type.my_resource resource_id
```

We also add an environment variable `TF_CLEAN_OUTPUT` which strips noisy Terraform output on a best effort basis. This is incompatible with interactive commands, so we only advise setting this in automation.


It additionally adds linters under the `lint` label for:
* `terraform fmt -check`
* `terraform validate`

See `//terraform/examples/<version>/BUILD` for examples of `terraform_root`. 

**NOTE**: This build rule utilises a [Terraform working directory](https://www.terraform.io/docs/cli/init/index.html) in `plz-out`, so whilst this is okay for demonstrations, you must use [Terraform Remote State](https://www.terraform.io/docs/language/state/remote.html) for your regular work. This can be added either simply through your `srcs` or through a `pre_binaries` binary.

---

## Usage

To use this build_def in your repository, you will need multiple files:
```
# advice: pick the latest commit
TERRAFORM_DEF_VERSION = "d0f1b02ae73893e695f1e50f9df2d6378f8701df"

# new Terraform rules
remote_file(
    name = "terraform",
    url = f"https://raw.githubusercontent.com/thought-machine/pleasings/{TERRAFORM_DEF_VERSION}/terraform/terraform.build_defs",
    hashes = ["95289dba7ae82131a7bb69976b5cdbedb4e7563c889a5b0d10da01d643be4540"],
    visibility = ["PUBLIC"],
)

# script for building Terraform modules
remote_file(
    name = "_terraform_module_builder",
    url = f"https://raw.githubusercontent.com/thought-machine/pleasings/{TERRAFORM_DEF_VERSION}/terraform/scripts/module_builder.sh",
    hashes = ["e1102ba9c15c29ebd98276912caf7bc45a7bc0a72780ccbbf4a077c0cc39b705"],
    visibility = ["PUBLIC"],
)

# script for running Terraform rules
remote_file(
    name = "_terraform_runner",
    url = f"https://raw.githubusercontent.com/thought-machine/pleasings/{TERRAFORM_DEF_VERSION}/terraform/scripts/runner.sh",
    hashes = ["144ca7da9037f07010547562f1e5c099811f1245c70e782e0e738c43a3072697"],
    visibility = ["PUBLIC"],
)

# script for building a Terraform workspace
remote_file(
    name = "_terraform_workspace_builder",
    url = f"https://raw.githubusercontent.com/thought-machine/pleasings/{TERRAFORM_DEF_VERSION}/terraform/scripts/workspace_builder.sh",
    hashes = ["ba435568d4ff9760aa72fc0bc8b92d07b606c7363f7bdefc5d4344f093d89590"],
    visibility = ["PUBLIC"],
)
```
