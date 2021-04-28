package(default_visibility = ["//visibility:public"])

load("@io_bazel_rules_go//go:def.bzl", "go_test", "go_library")

go_library(
  name = "go_default_library",
  srcs = glob(["*.go"]),
  importmap = "dev.io/cloud/vendor/dev.io/utils",
  importpath = "dev.io/cloud/utils",
  deps = [
    "@com_github_gorilla_mux//:go_default_library",
    "@org_golang_google_grpc//:go_default_library",
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
