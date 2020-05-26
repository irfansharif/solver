load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository", "new_git_repository")

local_repository(
    name = "ortools",
    path = "/home/red/git/or-tools-fork",
)

# git_repository(
#     name = "ortools",
#     commit = "5ff76b487a6c2006326765d6417964599eedc8c9",
#     remote = "https://github.com/google/or-tools.git",
# )

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

git_repository(
    name = "bazel_skylib",
    commit = "3721d32",  # release 0.8.0
    remote = "https://github.com/bazelbuild/bazel-skylib.git",
)

git_repository(
    name = "com_google_protobuf",
    commit = "fe1790c",  # release v3.11.2
    remote = "https://github.com/protocolbuffers/protobuf.git",
)

git_repository(
    name = "com_google_protobuf_cc",
    commit = "fe1790c",  # release v3.11.2
    remote = "https://github.com/protocolbuffers/protobuf.git",
)

load("@com_google_protobuf//:protobuf_deps.bzl", "protobuf_deps")

# Load common dependencies.
protobuf_deps()

git_repository(
    name = "com_google_absl",
    commit = "b832dce",  # release 20200225
    remote = "https://github.com/abseil/abseil-cpp.git",
)

http_archive(
    name = "gtest",
    build_file = "//bazel:gtest.BUILD",
    strip_prefix = "googletest-release-1.8.0/googletest",
    url = "https://github.com/google/googletest/archive/release-1.8.0.zip",
)

http_archive(
    name = "glpk",
    build_file = "//bazel:glpk.BUILD",
    sha256 = "9a5dab356268b4f177c33e00ddf8164496dc2434e83bd1114147024df983a3bb",
    url = "http://ftp.gnu.org/gnu/glpk/glpk-4.52.tar.gz",
)

# Go stuff
http_archive(
    name = "io_bazel_rules_go",
    urls = [
        "https://mirror.bazel.build/github.com/bazelbuild/rules_go/releases/download/v0.22.2/rules_go-v0.22.2.tar.gz",
        "https://github.com/bazelbuild/rules_go/releases/download/v0.22.2/rules_go-v0.22.2.tar.gz",
    ],
    sha256 = "142dd33e38b563605f0d20e89d9ef9eda0fc3cb539a14be1bdb1350de2eda659",
)

http_archive(
    name = "bazel_gazelle",
    urls = [
        "https://storage.googleapis.com/bazel-mirror/github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
        "https://github.com/bazelbuild/bazel-gazelle/releases/download/v0.20.0/bazel-gazelle-v0.20.0.tar.gz",
    ],
    sha256 = "d8c45ee70ec39a57e7a05e5027c32b1576cc7f16d9dd37135b0eddde45cf1b10",
)

load("@io_bazel_rules_go//go:deps.bzl", "go_rules_dependencies", "go_register_toolchains")

go_rules_dependencies()

go_register_toolchains()

load("@bazel_gazelle//:deps.bzl", "gazelle_dependencies", "go_repository")

gazelle_dependencies()

# git_repository(
#     name = "com_github_swig_swig",
#     commit = "8b572399d72f3d812165e0975498c930ae822a4f",
#     remote = "https://github.com/swig/swig.git",
# )

http_archive(
    name = "swig",
    build_file = "//third_party:swig.BUILD",
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
    build_file = "//third_party:pcre.BUILD",
    sha256 = "69acbc2fbdefb955d42a4c606dfde800c2885711d2979e356c0636efde9ec3b5",
    strip_prefix = "pcre-8.42",
    urls = [
        "https://storage.googleapis.com/mirror.tensorflow.org/ftp.exim.org/pub/pcre/pcre-8.42.tar.gz",
        "https://ftp.exim.org/pub/pcre/pcre-8.42.tar.gz",
    ],
)

go_repository(
    name = "com_github_google_go_cmp",
    importpath = "github.com/google/go-cmp",
    sum = "h1:/exdXoGamhu5ONeUJH0deniYLWYvQwW66yvlfiiKTu0=",
    version = "v0.4.1",
)

go_repository(
    name = "org_golang_x_xerrors",
    importpath = "golang.org/x/xerrors",
    sum = "h1:E7g+9GITq07hpfrRu66IVDexMakfv52eLZ2CXBWiKr4=",
    version = "v0.0.0-20191204190536-9bdfabe68543",
)