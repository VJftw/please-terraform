subinclude("//build/defs:terraform")

terraform_module(
    name = "my_module",
    srcs = [
        "data.tf",
        "main.tf",
    ],
    visibility = ["//example/0.11/..."],
    deps = [
        "//example/third_party/terraform/module:cloudposse_null_label_0_11",
    ],
)
