subinclude("///go//build_defs:go")

if "RELEASE" in CONFIG and CONFIG.RELEASE:
    architectures = [
        "darwin_amd64",
        "darwin_arm64",
        "freebsd_amd64",
        "linux_amd64",
        "linux_arm64",
        "windows_amd64",
    ]

    go_toolchain(
        name = "release_toolchain",
        architectures = architectures,
        version = "1.21.6",
    )

    def cross_compile(arch):
        suffix = arch
        if "windows" in arch:
            suffix += ".exe"

        return build_rule(
            name = f"please_terraform_{arch}",
            srcs = [f"///{arch}//cmd/please_terraform"],
            outs = [f"please_terraform_{suffix}"],
            cmd = "mv $SRC $OUT",
        )

    filegroup(
        name = f"release_files",
        srcs = [cross_compile(arch) for arch in architectures],
        labels = [
            f"hlink:plz-out/package",
        ],
    )
