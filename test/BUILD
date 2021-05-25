load("@io_bazel_rules_go//go:def.bzl", "go_test", "go_library")

go_test(
  name = "go_default_testing",
  srcs = glob(['*.go']),
  tags = ['selftest'],
  deps = [
    "//staging/src/dev.io/utils:go_default_library",
    "//staging/src/dev.io/protoc:go_default_library",
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
