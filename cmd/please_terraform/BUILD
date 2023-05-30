subinclude("///go//build_defs:go")

go_binary(
    name = "please_terraform",
    srcs = ["main.go"],
    static = True,
    visibility = ["PUBLIC"],
    deps = [
        "//internal/cmd",
        "//pkg/module",
        "//pkg/root",
    ],
)
