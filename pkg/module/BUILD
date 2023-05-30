subinclude("///go//build_defs:go")

go_library(
    name = "module",
    srcs = [
        "command.go",
        "local.go",
        "module.go",
        "registry.go",
    ],
    visibility = [
        "//cmd/...",
        "//pkg/...",
    ],
    deps = [
        "///third_party/go/github.com_hashicorp_go-getter//:go-getter",
        "//internal/logging",
        "//pkg/please",
    ],
)

go_test(
    name = "module_test",
    srcs = ["local_test.go"],
    external = True,
    deps = [
        ":module",
        "///third_party/go/github.com_stretchr_testify//assert",
        "///third_party/go/github.com_stretchr_testify//require",
    ],
)
