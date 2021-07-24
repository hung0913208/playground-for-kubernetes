load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

# Download gazelle to build go's dependencies
http_archive(
    name = "bazel_gazelle",
    sha256 = "222e49f034ca7a1d1231422cdb67066b885819885c356673cb1f72f748a3c9d4",
    urls = [
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.3/bazel-gazelle-v0.22.3.tar.gz",
    ],
)

# Download bazel rules to build docker image
http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "df13123c44b4a4ff2c2f337b906763879d94871d16411bf82dcfeba892b58607",
    strip_prefix = "rules_docker-0.13.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.13.0/rules_docker-v0.13.0.tar.gz"],
)

# Download bazel rules for nodejs
http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "dd4dc46066e2ce034cba0c81aa3e862b27e8e8d95871f567359f7a534cccb666",
    urls = [
        "https://github.com/bazelbuild/rules_nodejs/releases/download/3.1.0/rules_nodejs-3.1.0.tar.gz"
    ],
)

# Download bazel rules for golang
http_archive(
    name = "io_bazel_rules_go",
    sha256 = "b725e6497741d7fc2d55fcc29a276627d10e43fa5d0bb692692890ae30d98d00",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.24.3/rules_go-v0.24.3.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.24.3/rules_go-v0.24.3.tar.gz",
    ],
)

# Download bazel rules for grpc
http_archive(
    name = "rules_proto_grpc",
    sha256 = "d771584bbff98698e7cb3cb31c132ee206a972569f4dc8b65acbdd934d156b33",
    strip_prefix = "rules_proto_grpc-2.0.0",
    urls = ["https://github.com/rules-proto-grpc/rules_proto_grpc/archive/2.0.0.tar.gz"],
)

load("@rules_proto_grpc//:repositories.bzl", "rules_proto_grpc_toolchains", "rules_proto_grpc_repos")
rules_proto_grpc_toolchains()
rules_proto_grpc_repos()

load("@rules_proto_grpc//go:repositories.bzl", rules_proto_grpc_go_repos = "go_repos")
rules_proto_grpc_go_repos()

load("@rules_proto//proto:repositories.bzl", "rules_proto_dependencies", "rules_proto_toolchains")
rules_proto_dependencies()
rules_proto_toolchains()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
go_rules_dependencies()

load("//:deps.bzl", "go_dependencies")

# gazelle:repository_macro deps.bzl%go_dependencies
go_dependencies()
go_register_toolchains()

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

#node_repositories(package_json = ["//:package.json"])
#yarn_install(
#    name = "npm",
#    package_json = "//:package.json",
#    quiet = True,
#    yarn_lock = "//:yarn.lock",
#)
