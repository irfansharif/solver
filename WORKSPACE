load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository", "new_git_repository")

# TODO(irfansharif): Update to latest or-tools. Some of these dependencies
# (like glog/gflags) will no longer be needed.

git_repository(
    name = "ortools",
    commit = "b37d9c786b69128f3505f15beca09e89bf078a89",  # tag v9.0
    remote = "https://github.com/google/or-tools.git",
)

git_repository(
    name = "com_github_gflags_gflags",
    commit = "e171aa2",  # release v2.2.2
    remote = "https://github.com/gflags/gflags.git",
)

git_repository(
    name = "com_github_glog_glog",
    commit = "96a2f23",  # release v0.4.0
    remote = "https://github.com/google/glog.git",
)

http_archive(
    name = "bazel_skylib",
    sha256 = "1c531376ac7e5a180e0237938a2536de0c54d93f5c278634818e0efc952dd56c",
    urls = [
        "https://github.com/bazelbuild/bazel-skylib/releases/download/1.0.3/bazel-skylib-1.0.3.tar.gz",
        "https://mirror.bazel.build/github.com/bazelbuild/bazel-skylib/releases/download/1.0.3/bazel-skylib-1.0.3.tar.gz",
    ],
)

git_repository(
    name = "com_google_protobuf",
    commit = "436bd78",  # release v3.15.8
    remote = "https://github.com/protocolbuffers/protobuf.git",
)

git_repository(
    name = "com_google_absl",
    commit = "e1d388e", # release 20210324.1
    remote = "https://github.com/abseil/abseil-cpp.git",
)

http_archive(
    name = "rules_cc",
    strip_prefix = "rules_cc-262ebec3c2296296526740db4aefce68c80de7fa",
    urls = ["https://github.com/bazelbuild/rules_cc/archive/262ebec3c2296296526740db4aefce68c80de7fa.zip"],
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

# Load common dependencies.
protobuf_deps()

http_archive(
    name = "gtest",
    build_file = "@ortools//bazel:gtest.BUILD",
    strip_prefix = "googletest-release-1.8.0/googletest",
    url = "https://github.com/google/googletest/archive/release-1.8.0.zip",
)

http_archive(
    name = "glpk",
    build_file = "@ortools//bazel:glpk.BUILD",
    sha256 = "4281e29b628864dfe48d393a7bedd781e5b475387c20d8b0158f329994721a10",
    url = "http://ftp.gnu.org/gnu/glpk/glpk-4.65.tar.gz",
)

http_archive(
    name = "scip",
    build_file = "@ortools//bazel:scip.BUILD",
    patches = ["@ortools//bazel:scip.patch"],
    sha256 = "033bf240298d3a1c92e8ddb7b452190e0af15df2dad7d24d0572f10ae8eec5aa",
    url = "https://github.com/google/or-tools/releases/download/v7.7/scip-7.0.1.tgz",
)

http_archive(
    name = "bliss",
    build_file = "@ortools//bazel:bliss.BUILD",
    patches = ["@ortools//bazel:bliss-0.73.patch"],
    sha256 = "f57bf32804140cad58b1240b804e0dbd68f7e6bf67eba8e0c0fa3a62fd7f0f84",
    url = "http://www.tcs.hut.fi/Software/bliss/bliss-0.73.zip",
)

http_archive(
    name = "zlib",
    build_file = "@com_google_protobuf//:third_party/zlib.BUILD",
    sha256 = "c3e5e9fdd5004dcb542feda5ee4f0ff0744628baf8ed2dd5d66f8ca1197cb1a1",
    strip_prefix = "zlib-1.2.11",
    urls = [
        "https://mirror.bazel.build/zlib.net/zlib-1.2.11.tar.gz",
        "https://zlib.net/zlib-1.2.11.tar.gz",
    ],
)

http_archive(
    name = "io_bazel_rules_go",
    sha256 = "8e968b5fcea1d2d64071872b12737bbb5514524ee5f0a4f54f5920266c261acb",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.28.0/rules_go-v0.28.0.zip",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.28.0/rules_go-v0.28.0.zip",
    ],
)

http_archive(
    name = "swig",
    build_file = "//c-deps/external:swig.BUILD",
    sha256 = "58a475dbbd4a4d7075e5fe86d4e54c9edde39847cdb96a3053d87cb64a23a453",
    strip_prefix = "swig-3.0.8",
    urls = [
        "https://storage.googleapis.com/mirror.tensorflow.org/ufpr.dl.sourceforge.net/project/swig/swig/swig-3.0.8/swig-3.0.8.tar.gz",
        "https://ufpr.dl.sourceforge.net/project/swig/swig/swig-3.0.8/swig-3.0.8.tar.gz",
        "https://pilotfiber.dl.sourceforge.net/project/swig/swig/swig-3.0.8/swig-3.0.8.tar.gz",
    ],
)

http_archive(
    name = "pcre",
    build_file = "//c-deps/external:pcre.BUILD",
    sha256 = "69acbc2fbdefb955d42a4c606dfde800c2885711d2979e356c0636efde9ec3b5",
    strip_prefix = "pcre-8.42",
    urls = [
        "https://storage.googleapis.com/mirror.tensorflow.org/ftp.exim.org/pub/pcre/pcre-8.42.tar.gz",
        "https://ftp.exim.org/pub/pcre/pcre-8.42.tar.gz",
    ],
)

git_repository(
    name = "bazel_gazelle",
    commit = "d038863ba2e096792c6bb6afca31f6514f1aeecd",
    remote = "https://github.com/bazelbuild/bazel-gazelle",
)

# Load up our go dependencies (the ones listed under go.mod). The `DEPS.bzl`
# file is kept up to date using the `update-repos` Gazelle command (see
# README).
#
# gazelle:repository_macro DEPS.bzl%go_deps
load("//:DEPS.bzl", "go_deps")

# VERY IMPORTANT that we call into this function to prefer our pinned versions
# of the dependencies to any that might be pulled in via functions like
# `go_rules_dependencies`, `gazelle_dependencies`, etc.
go_deps()

load("@io_bazel_rules_go//go:deps.bzl", "go_register_toolchains", "go_rules_dependencies")

go_rules_dependencies()

go_register_toolchains(go_version = "1.16.5")

# Load gazelle dependencies.
load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies")

gazelle_dependencies()
