subinclude("//build/defs:terraform")

terraform_module(
    name = "my_module",
    srcs = ["main.tf"],
    visibility = ["//example/1.2/..."],
    deps = [
        "//example/third_party/terraform/module:cloudposse_null_label_0_12",
    ],
)
