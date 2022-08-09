package(default_visibility = ["PUBLIC"])

go_toolchain(
    name = "toolchain",
    hashes = [
        "9e5de37f9c49942c601b191ac5fba404b868bfc21d446d6960acc12283d6e5f2",  # linux-amd64
    ],
    version = "1.18.5",
)

go_module(
    name = "golint",
    binary = True,
    hashes = [
        "6a2982682499be8804589a0082514d650ca4019a88daeed951a215d56ff4f8b8",
    ],
    install = ["golint"],
    module = "golang.org/x/lint",
    version = "83fdc39ff7b56453e3793356bcff3070b9b96445",
    deps = [
        ":x_tools",
    ],
)

go_module(
    name = "x_tools",
    hashes = ["dff907b66b1f2cfd3536c92d71dfd5f28078026249f850214ee1a860519decc6"],
    install = [
        "cover",
        "present",
        "go/packages",
        "go/internal/packagesdriver",
        "internal/gocommand",
        "internal/packagesinternal",
        "go/gcexportdata",
        "go/internal/gcimporter",
        "go/ast/astutil",
    ],
    licences = ["bsd-3-clause"],
    module = "golang.org/x/tools",
    version = "5bcca83a78812bd91ce8cb29be4fc4521cdc8f6f",
    deps = [
        ":x_errors",
        ":x_mod",
    ],
)

go_module(
    name = "x_mod",
    hashes = ["2a937734dda11bfc94447ce43128e5ffafcf0479a4f4a4558aa17baf7335c2d7"],
    install = [
        "semver",
        "module",
    ],
    licences = ["bsd-3-clause"],
    module = "golang.org/x/mod",
    version = "v0.4.2",
    deps = [":x_errors"],
)

go_module(
    name = "x_errors",
    hashes = ["10780b586a60adfdea2f4d76dbaf3025eb3a8165376153c7febfb36516d1ec72"],
    install = [
        ".",
        "...",
    ],
    licences = ["bsd-3-clause"],
    module = "golang.org/x/xerrors",
    version = "a5947ffaace3e882f334c1750858b4a6a7e52422",
    deps = [":x_sys"],
)

go_module(
    name = "jessevdk_flags",
    module = "github.com/jessevdk/go-flags",
    version = "v1.5.0",
    deps = [
        ":x_sys",
    ],
)

go_module(
    name = "rs_zerolog",
    install = [
        ".",
        "internal/json",
    ],
    licences = ["MIT"],
    module = "github.com/rs/zerolog",
    version = "v1.23.0",
)

go_module(
    name = "hashicorp_getter",
    install = [
        ".",
        "helper/url",
    ],
    licences = ["MPL-2.0"],
    module = "github.com/hashicorp/go-getter",
    version = "v1.6.1",
    deps = [
        ":aws_aws_sdk_go",
        ":bgentry_netrc",
        ":cloud_google_storage",
        ":hashicorp_cleanhttp",
        ":hashicorp_safetemp",
        ":hashicorp_version",
        ":klauspost_compress",
        ":mitchellh_homedir",
        ":mitchellh_testing-interface",
        ":ulikunitz_xz",
    ],
)

go_module(
    name = "hashicorp_safetemp",
    licences = ["MPL-2.0"],
    module = "github.com/hashicorp/go-safetemp",
    version = "v1.0.0",
)

go_module(
    name = "mitchellh_testing-interface",
    licences = ["MIT"],
    module = "github.com/mitchellh/go-testing-interface",
    version = "v1.14.1",
)

go_module(
    name = "ulikunitz_xz",
    install = [
        ".",
        "lzma",
        "internal/...",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/ulikunitz/xz",
    version = "v0.5.10",
)

go_module(
    name = "klauspost_compress",
    install = [
        "...",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/klauspost/compress",
    version = "v1.12.2",
    deps = [":snappy"],
)

go_module(
    name = "snappy",
    licences = ["BSD-3-Clause"],
    module = "github.com/golang/snappy",
    version = "v0.0.3",
)

go_module(
    name = "hashicorp_cleanhttp",
    licences = ["MPL-2.0"],
    module = "github.com/hashicorp/go-cleanhttp",
    version = "v0.5.2",
)

go_module(
    name = "cloud_google_storage",
    licences = ["Apache-2.0"],
    module = "cloud.google.com/go/storage",
    version = "v1.15.0",
    deps = [
        ":cloud_google",
        ":google_api",
        ":google_genproto",
    ],
)

go_module(
    name = "cloud_google_compute",
    install = ["compute/metadata"],
    licences = ["Apache-2.0"],
    module = "cloud.google.com/go",
    version = "v0.82.0",
    deps = [
        ":google_genproto",
        ":google_grpc",
        ":opencensus",
    ],
)

go_module(
    name = "cloud_google",
    install = [
        "internal",
        "internal/version",
        "internal/trace",
        "internal/optional",
        "iam",
    ],
    licences = ["Apache-2.0"],
    module = "cloud.google.com/go",
    version = "v0.82.0",
    deps = [
        ":google_api",
        ":google_genproto",
        ":google_genproto_iam",
        ":google_grpc",
        ":opencensus",
    ],
)

go_module(
    name = "opencensus",
    install = [
        ".",
        "internal",
        "trace/internal",
        "trace",
        "trace/tracestate",
        "trace/propagation",
        "plugin/ochttp",
        "stats",
        "tag",
        "stats/view",
        "plugin/ochttp/propagation/b3",
        "metric/metricdata",
        "stats/internal",
        "internal/tagencoding",
        "resource",
        "metric/metricproducer",
    ],
    licences = ["Apache-2.0"],
    module = "go.opencensus.io",
    strip = ["examples"],
    version = "v0.23.0",
    deps = [":groupcache"],
)

go_module(
    name = "groupcache",
    install = ["lru"],
    licences = ["Apache-2.0"],
    module = "github.com/golang/groupcache",
    version = "v0.0.0-20210331224755-41bb18bfe9da",
)

go_module(
    name = "google_api",
    install = [
        "option",
        "option/internaloption",
        "googleapi",
        "storage/v1",
        "iterator",
        "transport/http",
        "internal/third_party/uritemplates",
        "internal/gensupport",
        "internal",
        "internal/impersonate",
        "googleapi/transport",
        "transport/cert",
        "transport/http/internal/propagation",
        "transport/internal/dca",
    ],
    licences = ["BSD-3-Clause"],
    module = "google.golang.org/api",
    version = "v0.47.0",
    deps = [
        ":googleapis_gax",
        ":x_oauth2",
    ],
)

go_module(
    name = "google_genproto_iam",
    install = [
        "googleapis/iam/v1",
    ],
    licences = ["Apache-2.0"],
    module = "google.golang.org/genproto",
    version = "v0.0.0-20210521181308-5ccab8a35a9a",
    deps = [
        ":google_grpc",
    ],
)

go_module(
    name = "google_genproto",
    install = [
        "googleapis/api/annotations",
        "googleapis/rpc/code",
        "googleapis/rpc/status",
        "googleapis/type/expr",
    ],
    licences = ["Apache-2.0"],
    module = "google.golang.org/genproto",
    version = "v0.0.0-20210521181308-5ccab8a35a9a",
    deps = [
        ":google_protobuf",
        ":protobuf",
    ],
)

go_module(
    name = "protobuf",
    install = [
        "proto",
        "ptypes",
        "ptypes/any",
        "ptypes/duration",
        "ptypes/timestamp",
    ],
    licences = ["BSD-3-Clause"],
    module = "github.com/golang/protobuf",
    version = "v1.5.2",
    deps = [":google_protobuf"],
)

go_module(
    name = "google_protobuf",
    install = [
        "proto",
        "reflect/protoreflect",
        "reflect/protoregistry",
        "reflect/protodesc",
        "runtime/protoimpl",
        "runtime/protoiface",
        "encoding/prototext",
        "encoding/protowire",
        "types/known/anypb",
        "types/known/durationpb",
        "types/known/timestamppb",
        "types/descriptorpb",
        "internal/pragma",
        "internal/filedesc",
        "internal/filetype",
        "internal/impl",
        "internal/version",
        "internal/encoding/messageset",
        "internal/encoding/text",
        "internal/errors",
        "internal/flags",
        "internal/genid",
        "internal/set",
        "internal/strs",
        "internal/order",
        "internal/encoding/defval",
        "internal/descfmt",
        "internal/descopts",
        "internal/encoding/tag",
        "internal/detrand",
    ],
    licences = ["BSD-3-Clause"],
    module = "google.golang.org/protobuf",
    version = "v1.26.0",
)

go_module(
    name = "google_grpc",
    install = [
        ".",
        "status",
        "codes",
        "metadata",
        "internal/status",
        "backoff",
        "balancer",
        "connectivity",
        "internal/buffer",
        "internal/channelz",
        "internal/grpcsync",
        "resolver",
        "balancer/base",
        "credentials",
        "internal/backoff",
        "internal/grpcutil",
        "internal/resolver",
        "internal/transport",
        "keepalive",
        "serviceconfig",
        "balancer/roundrobin",
        "internal/resolver/dns",
        "internal/resolver/passthrough",
        "internal/resolver/unix",
        "encoding",
        "encoding/proto",
        "internal",
        "internal/envconfig",
        "stats",
        "peer",
        "grpclog",
        "internal/binarylog",
        "internal/grpcrand",
        "tap",
        "internal/serviceconfig",
        "internal/balancerload",
        "attributes",
        "internal/credentials",
        "internal/metadata",
        "internal/syscall",
        "internal/transport/networktype",
        "balancer/grpclb/state",
        "internal/grpclog",
        "binarylog/grpc_binarylog_v1",
    ],
    licences = ["Apache-2.0"],
    module = "google.golang.org/grpc",
    version = "v1.38.0",
    deps = [
        ":google_genproto",
        ":protobuf",
        ":x_net",
        ":x_sys",
    ],
)

go_module(
    name = "googleapis_gax",
    licences = ["BSD-3-Clause"],
    module = "github.com/googleapis/gax-go/v2",
    version = "v2.0.5",
    deps = [":google_grpc"],
)

go_module(
    name = "x_sys",
    install = [
        "unix",
        "internal/...",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/sys",
    version = "v0.0.0-20210521203332-0cec03c779c1",
)

go_module(
    name = "x_net",
    install = [
        "http2",
        "http2/hpack",
        "trace",
        "context/ctxhttp",
        "http/httpguts",
        "idna",
        "internal/timeseries",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/net",
    version = "v0.0.0-20210521195947-fe42d452be8f",
    deps = [":x_text"],
)

go_module(
    name = "x_text",
    install = [
        "secure/bidirule",
        "unicode/bidi",
        "unicode/norm",
        "transform",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/text",
    version = "v0.3.6",
)

go_module(
    name = "x_oauth2",
    install = [
        ".",
        "google",
        "authhandler",
        "google/internal/externalaccount",
        "jwt",
        "internal",
        "jws",
    ],
    licences = ["BSD-3-Clause"],
    module = "golang.org/x/oauth2",
    version = "v0.0.0-20210514164344-f6687ab2804c",
    deps = [
        ":cloud_google_compute",
        ":x_net",
    ],
)

go_module(
    name = "hashicorp_version",
    licences = ["MPL-2.0"],
    module = "github.com/hashicorp/go-version",
    version = "v1.3.0",
)

go_module(
    name = "aws_aws_sdk_go",
    install = [
        "aws",
        "aws/credentials",
        "aws/credentials/ec2rolecreds",
        "aws/ec2metadata",
        "aws/session",
        "service/s3",
        "aws/endpoints",
        "aws/awserr",
        "internal/sdkio",
        "internal/sync/singleflight",
        "internal/ini",
        "internal/shareddefaults",
        "aws/client",
        "aws/request",
        "internal/sdkuri",
        "aws/client/metadata",
        "aws/corehandlers",
        "aws/credentials/processcreds",
        "aws/credentials/ssocreds",
        "aws/credentials/stscreds",
        "aws/defaults",
        "aws/csm",
        "aws/awsutil",
        "aws/signer/v4",
        "internal/s3shared/arn",
        "private/checksum",
        "private/protocol",
        "private/protocol/eventstream",
        "private/protocol/eventstream/eventstreamapi",
        "private/protocol/rest",
        "private/protocol/restxml",
        "internal/s3shared/s3err",
        "aws/arn",
        "internal/s3shared",
        "private/protocol/xml/xmlutil",
        "internal/sdkrand",
        "service/sso",
        "service/sso/ssoiface",
        "service/sts",
        "service/sts/stsiface",
        "aws/credentials/endpointcreds",
        "internal/strings",
        "internal/sdkmath",
        "private/protocol/query",
        "private/protocol/restjson",
        "private/protocol/json/jsonutil",
        "private/protocol/query/queryutil",
        "private/protocol/jsonrpc",
    ],
    licences = ["Apache-2.0"],
    module = "github.com/aws/aws-sdk-go",
    version = "v1.38.45",
    deps = [":jmespath_jmespath"],
)

go_module(
    name = "jmespath_jmespath",
    licences = ["Apache-2.0"],
    module = "github.com/jmespath/go-jmespath",
    version = "v0.4.0",
)

go_module(
    name = "bgentry_netrc",
    install = ["netrc"],
    licences = ["MIT"],
    module = "github.com/bgentry/go-netrc",
    version = "v0.0.0-20140422174119-9fd32a8b3d3d",
)

go_module(
    name = "mitchellh_homedir",
    licences = ["MIT"],
    module = "github.com/mitchellh/go-homedir",
    version = "v1.1.0",
)

go_module(
    name = "stretchr_testify",
    install = [
        "assert",
        "require",
    ],
    licences = ["MIT"],
    module = "github.com/stretchr/testify",
    version = "54d05a4e18440fb0a44989ab910f347eac4a131f",
    deps = [
        ":davecgh_spew",
        ":pmezard_difflib",
        ":stretchr_objx",
        ":yaml_v3",
    ],
)

go_module(
    name = "davecgh_spew",
    install = ["spew"],
    licences = ["ISC"],
    module = "github.com/davecgh/go-spew",
    version = "d8f796af33cc11cb798c1aaeb27a4ebc5099927d",
)

go_module(
    name = "pmezard_difflib",
    install = ["difflib"],
    module = "github.com/pmezard/go-difflib",
    version = "5d4384ee4fb2527b0a1256a821ebfc92f91efefc",
)

go_module(
    name = "stretchr_objx",
    licences = ["MIT"],
    module = "github.com/stretchr/objx",
    version = "35313a95ee26395aa17d366c71a2ccf788fa69b6",
    deps = [":davecgh_spew"],
)

go_module(
    name = "yaml_v3",
    licences = ["Apache-2.0"],
    module = "gopkg.in/yaml.v3",
    version = "496545a6307b2a7d7a710fd516e5e16e8ab62dbc",
)
