load("@io_bazel_rules_go//go:def.bzl", "go_library")
load("@rules_proto_grpc//go:defs.bzl", "go_grpc_library")

proto_library(
	name = "devio_proto",
	srcs = [
		"version.proto",
		"gateway.proto",
	],
	visibility = ["//visibility:public"],
)

go_grpc_library(
	name = "devio_go_proto",
	importpath = "dev.io/cloud/protoc",
  visibility = ["//visibility:public"],
	go_deps = [
		"@com_github_golang_protobuf//ptypes/any:go_default_library"
	],
	deps = [
		":devio_proto",
	]
)

go_library(
	name = "go_default_library",
	embed = [":devio_go_proto"],
	importpath = "dev.io/cloud/protoc",
	visibility = ["//visibility:public"],
)

