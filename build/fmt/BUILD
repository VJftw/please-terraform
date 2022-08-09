sh_cmd(
    name = "go",
    srcs = ["go.sh"],
    cmd = """
source $(out_location //build/util)
export GO_ROOT=$(out_location //third_party/go:toolchain)
export GO_FMT="\\\\$GO_ROOT/bin/gofmt"
export GO_BIN="\\\\$GO_ROOT/bin/go"
source $SRCS
    """,
    data = [
        "//build/util",
        "//third_party/go:toolchain|go",
    ],
    labels = ["fmt"],
    shell = "/bin/bash",
)

sh_cmd(
    name = "plz",
    srcs = ["plz.sh"],
    cmd = """
source $(out_location //build/util)
source $SRCS
    """,
    data = [
        "//build/util",
    ],
    labels = ["fmt"],
    shell = "/bin/bash",
)
