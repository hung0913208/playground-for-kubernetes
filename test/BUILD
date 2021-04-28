package(default_visibility = ["//visibility:public"])

filegroup(
	name = "package-srcs",
	srcs = glob(["**"]),
	visibility = ["//visibility:private"]
)

filegroup(
	name = "all-srcs",
	srcs = [
		":package-srcs",
		"//test/multiplexer:all-srcs",
	]
)
