sh_cmd(
    name = "go",
    srcs = ["go.sh"],
    cmd = """
source $(out_location //build/util)
export GO_LINT=$(out_location //third_party/go:golint)
source $SRCS
    """,
    data = [
        "//build/util",
        "//third_party/go:golint",
    ],
    labels = ["lint"],
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
    labels = ["lint"],
    shell = "/bin/bash",
)
