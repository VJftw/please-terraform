subinclude("///go//build_defs:go")

go_library(
    name = "logging",
    srcs = ["logging.go"],
    visibility = [
        "//cmd/...",
        "//internal/...",
        "//pkg/...",
    ],
    deps = [
        "///third_party/go/github.com_rs_zerolog//:zerolog",
    ],
)
