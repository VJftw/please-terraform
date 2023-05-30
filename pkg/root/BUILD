subinclude("///go//build_defs:go")

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
        "///third_party/go/github.com_stretchr_testify//assert",
    ],
)
