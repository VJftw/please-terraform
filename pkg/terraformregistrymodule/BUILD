go_library(
    name = "terraformregistrymodule",
    srcs = [
        "build.go",
        "terraformregistrymodule.go",
    ],
    visibility = ["//cmd/..."],
    deps = [
        "//internal/logging",
        "//pkg/plz",
        "//pkg/terraformregistryprovider",
        "//third_party/go:hashicorp_getter",
        "//third_party/go:hashicorp_hcl_v2",
    ],
)
