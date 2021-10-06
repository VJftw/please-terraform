go_library(
    name = "module",
    srcs = [
        "local.go",
        "module.go",
        "registry.go",
    ],
    visibility = [
        "//cmd/...",
        "//pkg/...",
    ],
    deps = [
        "//internal/logging",
        "//pkg/please",
        "//third_party/go:hashicorp_getter",
    ],
)

go_test(
    name = "module_test",
    srcs = ["module_test.go"],
    data = [
        # "testdata/main.tf",
    ],
    external = True,
    deps = [
        ":module",
        "//third_party/go:stretchr_testify",
    ],
)
