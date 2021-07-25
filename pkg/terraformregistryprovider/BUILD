go_library(
    name = "terraformregistryprovider",
    srcs = [
        "build.go",
        "terraformregistryprovider.go",
    ],
    visibility = [
        "//cmd/...",
        "//pkg/...",
    ],
    deps = [
        "//pkg/plz",
        "//third_party/go:hashicorp_getter",
        "//third_party/go:hashicorp_hcl_v1",
        "//third_party/go:hashicorp_hcl_v2",
    ],
)
