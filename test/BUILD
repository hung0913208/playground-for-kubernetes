load("@io_bazel_rules_go//go:def.bzl", "go_test", "go_library")

go_test(
  name = "test_muxer",
  srcs = [
    "muxer.go",
  ],
  tags = ["selftest"],
  deps = [
    "//staging/src/dev.io/utils:go_default_library",
  ]
)

go_test(
  name = "test_rpc",
  srcs = [
    "rpc.go",
  ],
  tags = ["selftest"],
  deps = [
    "//staging/src/dev.io/utils:go_default_library",
    "//staging/src/dev.io/protoc:go_default_library",
    "@org_golang_google_grpc//:go_default_library",
  ]
)

filegroup(
	name = "package-srcs",
	srcs = glob(["**"]),
	visibility = ["//visibility:private"]
)

filegroup(
	name = "all-srcs",
	srcs = [
		":package-srcs",
	]
)
