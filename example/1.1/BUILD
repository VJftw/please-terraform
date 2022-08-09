subinclude("//build/defs:terraform")

terraform_root(
    name = "example",
    srcs = [
        "main.tf",
        "variables.tf",
    ],
    modules = [
        "//example/third_party/terraform/module:cloudposse_null_label_0_12",
        "//example/1.1/my_module:my_module",
    ],
    toolchain = "//example/third_party/terraform:1.1",
    var_files = ["my_vars.tfvars"],
)
