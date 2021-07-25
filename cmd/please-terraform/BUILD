go_binary(
    name = "please-terraform",
    srcs = ["main.go"],
    visibility = ["PUBLIC"],
    deps = [
        "//pkg/terraformregistrymodule",
        # "//pkg/terraformmodule",
        "//pkg/terraformregistryprovider",
        "//third_party/go:jessevdk_flags",
        "//internal/cmd",
    ],
)
