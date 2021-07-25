load("@rules_proto_grpc//go:defs.bzl", "go_grpc_library")
load("@io_bazel_rules_go//go:def.bzl", "go_library")

proto_library(
	name = "gateway_proto",
	srcs = [
		"gateway.proto",
	],
	visibility = ["//visibility:public"],
	deps = [
    "@com_github_golang_protobuf//:any_proto",  # Well-known	
	]
)

go_grpc_library(
	name = "gateway_go_proto",
	importpath = "dev.io/cloud/protoc",
  visibility = ["//visibility:public"],
	proto = ":gateway_proto",
)

go_library(
	name = "go_default_library",
	deps = [
		":gateway_go_proto",
	],
	importpath = "dev.io/cloud/protoc",
	visibility = ["//visibility:public"],
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
