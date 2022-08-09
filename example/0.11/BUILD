subinclude("//build/defs:terraform")

terraform_root(
    name = "0.11",
    srcs = [
        "data.tf",
        "main.tf",
        "variables.tf",
    ],
    modules = [
        "//example/third_party/terraform/module:cloudposse_null_label_0_11",
        "//example/0.11/my_module:my_module",
    ],
    toolchain = "//example/third_party/terraform:0.11",
    var_files = ["my_vars.tfvars"],
)
