subinclude("//build/defs:terraform")

terraform_root(
    name = "example",
    srcs = [
        "data.tf",
        "main.tf",
        "variables.tf",
    ],
    modules = [
        "//example/third_party/terraform/module:cloudposse_null_label_0_11",
        "//example/0.11/my_module:my_module",
    ],
    providers = [
        "//example/third_party/terraform/provider:null",
    ],
    toolchain = "//example/third_party/terraform:terraform_0_11",
    var_files = ["my_vars.tfvars"],
)
