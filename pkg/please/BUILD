subinclude("///go//build_defs:go")

go_library(
    name = "please",
    srcs = [
        "please.go",
        "replace.go",
        "sync.go",
    ],
    visibility = ["//pkg/..."],
    deps = [
        "//internal/logging",
    ],
)

go_test(
    name = "please_test",
    srcs = [
        "sync_test.go",
    ],
    external = True,
    deps = [
        ":please",
        "///third_party/go/github.com_stretchr_testify//assert",
        "///third_party/go/github.com_stretchr_testify//require",
    ],
)
