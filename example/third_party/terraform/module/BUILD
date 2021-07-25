subinclude("//build/defs:terraform")

terraform_registry_module(
    name = "cloudposse_null_label_0_11",
    # hashes = ["519b74c3ae8da07da1816499f502926c35a1d17914b5a814dfb81c2c41a59b75"],
    licences = ["Apache-2.0"],
    module = "cloudposse/label/null",
    strip = ["examples"],
    version = "0.11.1",
    visibility = ["PUBLIC"],
)

terraform_registry_module(
    name = "cloudposse_null_label_0_12",
    # hashes = ["699153c70268b51ebbb826a9f5cf97ae5e12576810a25d8de421f1b6562f78e9"],
    licences = ["Apache-2.0"],
    module = "cloudposse/label/null",
    strip = [
        "examples",
        "exports",
    ],
    version = "0.22.1",
    visibility = ["PUBLIC"],
)

terraform_registry_module(
    name = "cloudposse_route53_cluster_hostname_0_12",
    # hashes = ["fe0e24ab7d161c582cd575cd34202e5ce3213f292d1b329a9523dbd5a085388c"],
    licences = ["Apache-2.0"],
    module = "cloudposse/route53-cluster-hostname/aws",
    strip = ["examples"],
    version = "0.12.0",
    visibility = ["PUBLIC"],
    deps = [
        ":cloudposse_null_label_0_12",
        "//example/third_party/terraform/provider:aws",
        "//example/third_party/terraform/provider:local",
    ],
)

terraform_registry_module(
    name = "aws_security_group",
    module = "terraform-aws-modules/security-group/aws",
    strip = ["examples"],
    version = "4.0.0",
    deps = [
        "//example/third_party/terraform/provider:aws",
    ],
)
