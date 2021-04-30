# gazelle:prefix
# gazelle:importmap_prefix dev.io/cloud

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
    "//staging/src/dev.io/ctrl:all-srcs",
    "//staging/src/dev.io/gw:all-srcs",
    "//staging/src/dev.io/utils:all-srcs",
  ],
  tags = ["automanaged"],
  visibility = ["//visibility:public"],
)
