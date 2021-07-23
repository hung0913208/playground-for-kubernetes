load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Download gazelle to build go's dependencies
http_archive(
    name = "bazel_gazelle",
    sha256 = "112ceace31ac48a9dde28f1f1ad98e76fc7f901ef088b944a84e55bc93cd198a",
    urls = [
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.3/bazel-gazelle-v0.22.3.tar.gz",
    ],
)

# Download bazel_rules to build docker image
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "df13123c44b4a4ff2c2f337b906763879d94871d16411bf82dcfeba892b58607",
    strip_prefix = "rules_docker-0.13.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.13.0/rules_docker-v0.13.0.tar.gz"],
)

# Download bazel_rules for nodejs
http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "bfacf15161d96a6a39510e7b3d3b522cf61cb8b82a31e79400a84c5abcab5347",
    urls = [
        "https://github.com/bazelbuild/rules_nodejs/releases/download/3.2.1/rules_nodejs-3.2.1.tar.gz"
    ],
)

# Download bazel rules for golang
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "aa301ab560203bf740d07456a505730bf1ee20f4c471f77357cd31e7e11f5170",
    urls = [
        "https://github.com/bazelbuild/rules_go/releases/download/v0.25.1/rules_go-v0.25.1.tar.gz",
    ],
)

# Download bazel rules for grpc
http_archive(
    name = "rules_proto_grpc",
    sha256 = "7954abbb6898830cd10ac9714fbcacf092299fda00ed2baf781172f545120419",
    strip_prefix = "rules_proto_grpc-3.1.1",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/3.1.1.tar.gz"],
)

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_toolchains", "rules_proto_grpc_repos")
rules_proto_grpc_toolchains()
rules_proto_grpc_repos()

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
rules_proto_dependencies()
rules_proto_toolchains()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
go_rules_dependencies()

load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()
go_register_toolchains(version = "1.15.8")

load("@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",)
load("@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",)

load("@build_bazel_rules_nodejs//:index.bzl", "node_repositories", "yarn_install")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

protobuf_deps()
gazelle_dependencies()

container_repositories()
_go_image_repos()

load("@rules_proto_grpc//go:repositories.bzl", rules_proto_grpc_go_repos = "go_repos")

rules_proto_grpc_go_repos()

#node_repositories(package_json = ["//:package.json"])
#yarn_install(
#    name = "npm",
#    package_json = "//:package.json",
#    quiet = True,
#    yarn_lock = "//:yarn.lock",
#)
