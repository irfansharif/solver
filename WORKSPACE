load("@bazel_tools//tools/build_defs/repo:http.bzl", "http_archive")
load("@bazel_tools//tools/build_defs/repo:git.bzl", "git_repository", "new_git_repository")

git_repository(
    name = "ortools",
    commit = "5ff76b487a6c2006326765d6417964599eedc8c9",
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
    commit = "b832dce", # release 20200225
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
