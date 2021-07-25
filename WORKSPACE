workspace(name = "cluster")

load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Download bazel rules to build docker image
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "df13123c44b4a4ff2c2f337b906763879d94871d16411bf82dcfeba892b58607",
    strip_prefix = "rules_docker-0.13.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.13.0/rules_docker-v0.13.0.tar.gz"],
)

http_archive(
    name = "rules_proto_grpc",
   sha256 = "7954abbb6898830cd10ac9714fbcacf092299fda00ed2baf781172f545120419",
    strip_prefix = "rules_proto_grpc-3.1.1",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/3.1.1.tar.gz"],
)

load("@rules_proto_grpc//:repositories.bzl",
     "rules_proto_grpc_toolchains",
     "rules_proto_grpc_repos",
     io_bazel_rules_go_repos = "io_bazel_rules_go",
     bazel_gazelle_repos = "bazel_gazelle")
rules_proto_grpc_toolchains()
rules_proto_grpc_repos()
io_bazel_rules_go_repos()
bazel_gazelle_repos()

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
rules_proto_dependencies()
rules_proto_toolchains()

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")
protobuf_deps()

load("@rules_proto_grpc//go:repositories.bzl", rules_proto_grpc_go_repos = "go_repos")
rules_proto_grpc_go_repos()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
go_rules_dependencies()

load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()
go_register_toolchains(version = "1.15.5")

load("@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",)
container_repositories()

load("@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",)
_go_image_repos()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
gazelle_dependencies()

#load("@build_bazel_rules_nodejs//:index.bzl", "node_repositories", "yarn_install")
#node_repositories(package_json = ["//:package.json"])
#yarn_install(
#    name = "npm",
#    package_json = "//:package.json",
#    quiet = True,
#    yarn_lock = "//:yarn.lock",
#)
