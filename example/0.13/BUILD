subinclude("//build/defs:terraform")

terraform_root(
    name = "example",
    srcs = [
        "main.tf",
        "variables.tf",
    ],
    modules = [
        "//example/third_party/terraform/module:cloudposse_null_label_0_12",
        "//example/0.13/my_module:my_module",
    ],
    providers = [
        "//example/third_party/terraform/provider:null",
    ],
    toolchain = "//example/third_party/terraform:terraform_0_13",
    var_files = ["my_vars.tfvars"],
)
