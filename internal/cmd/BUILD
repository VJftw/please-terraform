subinclude("///go//build_defs:go")

go_library(
    name = "cmd",
    srcs = [
        "flags.go",
        "logging.go",
    ],
    visibility = [
        "//build/...",
        "//cmd/...",
    ],
    deps = [
        "///third_party/go/github.com_jessevdk_go-flags//:go-flags",
        "///third_party/go/github.com_rs_zerolog//:zerolog",
        "//internal/logging",
    ],
)
