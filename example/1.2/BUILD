subinclude("//build/defs:terraform")

terraform_root(
    name = "example",
    srcs = [
        "main.tf",
        "variables.tf",
    ],
    modules = [
        "//example/1.2/my_module:my_module",
    ],
    toolchain = "//example/third_party/terraform:1.2",
    var_files = ["my_vars.tfvars"],
)
