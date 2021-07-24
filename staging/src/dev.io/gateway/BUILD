package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_test", "go_library")

go_library(
  name = "go_default_library",
  srcs = glob(["*.go"]),
  importmap = "dev.io/cloud/vendor/dev.io/gw",
  importpath = "dev.io/cloud/gw",
  deps = [
    "//staging/src/dev.io/protoc:go_default_library",
    "//staging/src/dev.io/utils:go_default_library",
    "@com_github_gorilla_mux//:go_default_library",
    "@org_golang_google_grpc//:go_default_library",
    "@com_github_google_uuid//:go_default_library",
    "@com_github_graphql-go_graphql//:go_default_library",
  ]
)

filegroup(
  name = "package-srcs",
  srcs = glob(["**"]),
  tags = ["automanaged"],
  visibility = ["//visibility:private"],
)

filegroup(
  name = "all-srcs",
  srcs = [
    ":package-srcs",
  ],
  tags = ["automanaged"],
  visibility = ["//visibility:public"],
)
