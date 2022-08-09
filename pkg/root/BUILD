go_library(
    name = "root",
    srcs = [
        "build.go",
        "command.go",
        "virtualenv.go",
    ],
    visibility = [
        "//cmd/...",
    ],
    deps = [
        "//internal/logging",
        "//pkg/module",
        "//pkg/please",
    ],
)

go_test(
    name = "root_test",
    srcs = [
        "build_test.go",
    ],
    external = True,
    deps = [
        ":root",
        "//third_party/go:stretchr_testify",
    ],
)
