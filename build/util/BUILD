subinclude("///shell//build_defs:shell")

sh_cmd(
    name = "util",
    srcs = ["util.sh"],
    cmd = """
source "$(out_location //third_party/sh:ansi)"
source $SRCS
    """,
    data = [
        "//third_party/sh:ansi",
    ],
    visibility = [
        "//build/...",
        "//test/...",
    ],
)
