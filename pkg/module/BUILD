go_library(
    name = "module",
    srcs = [
        "module.go",
        "local.go",
        "registry.go",
    ],
    visibility = ["//cmd/...", "//pkg/..."],
    deps = [
        "//internal/logging",
        "//pkg/please",
         "//third_party/go:hashicorp_getter",
    ]
)

go_test(
    name = "module_test",
    srcs = ["module_test.go"],
    deps = [
        ":module",
        "//third_party/go:stretchr_testify",
    ],
    data = [
        # "testdata/main.tf",
    ],
    external = True,
)
