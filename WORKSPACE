load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")

http_archive(
  name = "io_bazel_rules_go",
  sha256 = "7904dbecbaffd068651916dce77ff3437679f9d20e1a7956bff43826e7645fcc",
  urls = [
    "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.25.1/rules_go-v0.25.1.tar.gz",
    "https://github.com/bazelbuild/rules_go/releases/download/v0.25.1/rules_go-v0.25.1.tar.gz",
  ],
)

http_archive(
  name = "bazel_gazelle",
  sha256 = "222e49f034ca7a1d1231422cdb67066b885819885c356673cb1f72f748a3c9d4",
  urls = [
    "https://mirror.bazel.build/github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.3/bazel-gazelle-v0.22.3.tar.gz",
    "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.22.3/bazel-gazelle-v0.22.3.tar.gz",
    ],
)

http_archive(
  name = "com_google_protobuf",
  sha256 = "9748c0d90e54ea09e5e75fb7fac16edce15d2028d4356f32211cfa3c0e956564",
  strip_prefix = "protobuf-3.11.4",
  urls = ["https://github.com/protocolbuffers/protobuf/archive/v3.11.4.zip"],
)

http_archive(
    name = "build_bazel_rules_typescript",
    sha256 = "136ba6be39b4ff934cc0f41f043912305e98cb62254d9e6af467e247daafcd34",
    strip_prefix = "rules_typescript-0.22.0",
    url = "https://github.com/bazelbuild/rules_typescript/archive/0.22.0.zip",
)

http_archive(
    name = "io_bazel_rules_docker",
    sha256 = "df13123c44b4a4ff2c2f337b906763879d94871d16411bf82dcfeba892b58607",
    strip_prefix = "rules_docker-0.13.0",
    urls = ["https://github.com/bazelbuild/rules_docker/releases/download/v0.13.0/rules_docker-v0.13.0.tar.gz"],
)

http_archive(
    name = "build_bazel_rules_nodejs",
    sha256 = "55a25a762fcf9c9b88ab54436581e671bc9f4f523cb5a1bd32459ebec7be68a8",
    urls = ["https://github.com/bazelbuild/rules_nodejs/releases/download/3.2.2/rules_nodejs-3.2.2.tar.gz"],
)

http_archive(
    name = "bazel_skylib",
    urls = [
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.0.3/bazel-skylib-1.0.3.tar.gz",
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.0.3/bazel-skylib-1.0.3.tar.gz",
    ],
    sha256 = "1c531376ac7e5a180e0237938a2536de0c54d93f5c278634818e0efc952dd56c",
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")
load("@build_bazel_rules_typescript//:package.bzl", "rules_typescript_dependencies")
load("@build_bazel_rules_nodejs//:index.bzl", "node_repositories", "yarn_install")
load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")
load("@io_bazel_rules_docker//repositories:repositories.bzl",
    container_repositories = "repositories",)
load("@io_bazel_rules_docker//go:image.bzl",
    _go_image_repos = "repositories",)

protobuf_deps()

go_rules_dependencies()
go_register_toolchains(version = "1.16")
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
#rules_typescript_dependencies()
