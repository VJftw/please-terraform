go_library(
    name = "provider",
    srcs = [
        "provider.go",
        "registry.go",
    ],
    visibility = [
        "//cmd/...",
        "//pkg/module/...",
    ],
    deps = [
        "//internal/logging",
        "//pkg/please",
        "//third_party/go:hashicorp_getter",
    ],
)
