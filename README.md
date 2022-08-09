# please-terraform
Terraform integration w/ the Please build system.

This includes support for the following:
 * `terraform_toolchain`: Easy management of multiple versions of Terraform.
 * `terraform_module`: Terraform modules from the local filesystem.
 * `terraform_registry_module`: Terraform Modules from the Terraform Registry.
 * `terraform_root`: Terraform root configuration management.


## `terraform_toolchain`

This build rule allows you to specify a Terraform version to download and re-use in `terraform_root` rules. You can repeat this for multiple versions if you like, see `//third_party/terraform/BUILD` for examples.

## `terraform_module`

This build rule allows you to specify a [Terraform module](https://www.terraform.io/docs/language/modules/index.html) to re-use in your `terraform_root` rules or as dependencies in other `terraform_module` rules. These Terraform modules are sourced locally on the filesystem. For externally sourced modules, see `terraform_registry_module`.

See `//examples/<version>/my_module/BUILD` for examples of `terraform_module`s.

In your Terraform source code, you should refer to your modules by their canonical build label. e.g.:

```typescript
module "my_module" {
    source = "//examples/0.12/my_module:my_module"
}
```

If your module has providers or required providers configuration, you must include them as deps.


## `terraform_registry_module`

This build rule allows you to specify a [Terraform Module]() from the a Terraform Registry to re-use in your `terraform_root` rules or as dependencies in other `terraform_module` rules.

Set `//example/third_party/terraform/module/BUILD` for examples of `terraform_registry_module`s.

```typescript
module "remote_module" {
    source = "//examples/third_party/modules:my_module"
}
```

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
 * `<name>`: for all workflows. This sets up a Virtual Environment where `terraform` can be called directly. For example:
    * `plz run //my_infrastructure_tf -- terraform init`
    * `plz run //my_infrastructure_tf -- "terraform init && terraform console"`
 * `_plan`
 * `_apply`
 * `_destroy`
 * `` for all other workflows e.g.

For all of these workflows, we support passing in flags via please as expected, e.g.:
```
$ plz run //my_tf:my_tf_plan -- -lock=false
$ plz run //my_tf:my_tf_import -- resource_type.my_resource resource_id
```

See `//example/<version>/BUILD` for examples of `terraform_root`.

**NOTE**: This build rule utilises a [Terraform working directory](https://www.terraform.io/docs/cli/init/index.html) in `plz-out`, so whilst this is okay for demonstrations, you must use [Terraform Remote State](https://www.terraform.io/docs/language/state/remote.html) for your regular work. This can be added either simply through your `srcs` or through a `pre_binaries` binary.

---

## Usage


### Please Plugin

```ini
# .plzconfig

## Support the please-* format of Please plugins.
PluginRepo = ["https://github.com/{owner}/{plugin}/archive/{revision}.zip"]

[Plugin "terraform"]
Target = //third_party/plugins:terraform
Tool = //third_party/plugins:terraform_tool
```

```python
# //third_party/plugins/BUILD
plugin_repo(
    name = "terraform",
    owner = "VJftw",
    plugin = "please-terraform",
    revision = "<version>",
)

remote_file(
    name = "terraform_tool",
    url = f"https://github.com/VJftw/please-terraform/releases/download/<version>/please-terraform"
    visibility = ["PUBLIC"],
    binary = True,
)
```

### Please Remote Files

```python
# //third_party/defs/BUILD
TERRAFORM_DEF_VERSION="<version>"
TERRAFORM_TOOL="//third_party/defs:terraform_tool"
TERRAFORM_DEFAULT_TOOLCHAIN="//third_party/terraform:1.0"

remote_file(
    name = "_terraform#download",
    url = f"https://raw.githubusercontent.com/VJftw/please-terraform/{TERRAFORM_DEF_VERSION}/build/defs/terraform.build_defs",
    hashes = ["95289dba7ae82131a7bb69976b5cdbedb4e7563c889a5b0d10da01d643be4540"],
)

remote_file(
    name = "terraform_tool",
    url = f"https://github.com/VJftw/please-terraform/releases/download/{TERRAFORM_DEF_VERSION}/please-terraform",
    visibility = ["PUBLIC"],
    binary = True,
)

genrule(
    name = "terraform",
    srcs = [":_terraform#download"],
    outs = ["terraform_custom.build_defs"],
    cmd = [
        "mv $SRCS $OUTS",
        # Replace CONFIG.TERRAFORM.TOOL with your tool.
        f"sed -i 's#CONFIG.TERRAFORM.TOOL#{TERRAFORM_TOOL}#g' $OUTS",
        # Replace CONFIG.TERRAFORM.DEFAULT_TOOLCHAIN with your tool.
        f"sed -i 's#CONFIG.TERRAFORM.DEFAULT_TOOLCHAIN#{TERRAFORM_DEFAULT_TOOLCHAIN}#g' $OUTS",
    ],
    visibility = ["PUBLIC"],
)
```


## Future Work

* `terraform_module`:
    * Support uploading these Terraform Modules to the Terraform Registry for promotion-based configuration.

* `terraform_root`:
    * Add optional additional rules for linting:
        * `terraform fmt -check`
        * `terraform init -lock=false && terraform validate`


### Future Work - Examples

- Extending with [Terratest](https://terratest.gruntwork.io/).
- Extending with [OPA](https://www.openpolicyagent.org/docs/latest/terraform/).
