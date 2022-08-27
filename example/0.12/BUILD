subinclude("//build/defs:terraform")

terraform_root(
    name = "example",
    srcs = [
        "main.tf",
        "outputs.tf",
        "variables.tf",
    ],
    modules = [
        "//example/third_party/terraform/module:cloudposse_null_label_0_12",
        "//example/0.12/my_module:my_module",
    ],
    toolchain = "//example/third_party/terraform:0.12",
    var_files = ["my_vars.tfvars"],
    vars = {
        "single_value": "e",
        "list_value": [
            "e_1",
            "e_2",
        ],
        "dict_value": {
            "e_1": "e_1",
            "e_2": "e_2",
        },
    },
)
