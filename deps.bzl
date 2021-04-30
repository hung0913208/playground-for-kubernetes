
load(
    "@bazel_gazelle//internal:go_repository.bzl",
    _go_repository = "go_repository",
)

load(
    "@bazel_gazelle//internal:go_repository_cache.bzl",
    "go_repository_cache",
)

load(
    "@bazel_gazelle//internal:go_repository_tools.bzl",
    "go_repository_tools",
)

load(
    "@bazel_gazelle//internal:go_repository_config.bzl",
    "go_repository_config",
)

load(
    "@bazel_tools//tools/build_defs/repo:git.bzl",
    _tools_git_repository = "git_repository",
)

def _maybe(repo_rule, name, **kwargs):
    if name not in native.existing_rules():
        repo_rule(name = name, **kwargs)

def _go_default_dependencies():
    pass

def go_dependencies():
    _go_default_dependencies()
