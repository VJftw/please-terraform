subinclude("//build/defs:terraform")

for version, hashes in {
    "0.11.15": ["e6c8c884de6c353cf98252c5e11faf972d4b30b5d070ab5ff70eaf92660a5aac"],
    "0.12.31": ["e5eeba803bc7d8d0cae7ef04ba7c3541c0abd8f9e934a5e3297bf738b31c5c6d"],
    "0.13.7": ["4a52886e019b4fdad2439da5ff43388bbcc6cce9784fde32c53dcd0e28ca9957"],
    "0.14.11": ["171ef5a4691b6f86eab524feaf9a52d5221c875478bd63dd7e55fef3939f7fd4"],
    "0.15.5": ["3b144499e08c245a8039027eb2b84c0495e119f57d79e8fb605864bb48897a7d"],
    "1.0.11": ["eeb46091a42dc303c3a3c300640c7774ab25cbee5083dafa5fd83b54c8aca664"],
    "1.1.9": ["9d2d8a89f5cc8bc1c06cb6f34ce76ec4b99184b07eb776f8b39183b513d7798a"],
    "1.2.6": ["9fd445e7a191317dcfc99d012ab632f2cc01f12af14a44dfbaba82e0f9680365"],
}.items():
    version_parts = version.split(".")
    major_version = version_parts[0]
    minor_version = version_parts[1]

    terraform_toolchain(
        name = f"{major_version}.{minor_version}",
        hashes = hashes,
        version = version,
        visibility = [f"//example/{major_version}.{minor_version}/..."],
    )
