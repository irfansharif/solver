help:
	@echo "Supported commands: build, test, generate, rewrite"

build:
	@bazel build ...
	@echo "ok"

test:
	@bazel test --test_output=all --test_arg='-test.v' ...
	@echo "ok"

dd = $(shell bazel query "rdeps(..., @com_github_cockroachdb_datadriven//:datadriven)")
rewrite:
	@bazel test $(dd) --test_filter='TestDatadriven' --test_output=all --test_arg='-test.v' --test_arg='-rewrite'
	@echo "ok"

swig: internal/*.i
	@echo "--- generating swig files"
	@swig -v -go -cgo -c++ -intgosize 64 \
		-Ic-deps/or-tools \
		-Ic-deps/abseil-cpp \
		-o internal/wrapper.cc \
		-module internal \
		internal/sat.i > /dev/null
	@gofmt -s -w internal/internal.go

proto: internal/pb/*.proto
	@echo "--- generating proto files"
	@protoc --proto_path=internal/pb \
		--go_out=internal/pb \
		--go_opt=Mcp_model.proto=github.com/irfansharif/solver/internal/pb \
		--go_opt=Msat_parameters.proto=github.com/irfansharif/solver/internal/pb \
		--go_opt=paths=source_relative \
		cp_model.proto sat_parameters.proto > /dev/null

bzl: FORCE
	@echo "--- generating bazel files"
	@echo "build --cxxopt=-std=c++17 --experimental_convenience_symlinks=ignore --define gotags=bazel" > $@.tmp
	@echo "test  --cxxopt=-std=c++17 --experimental_convenience_symlinks=ignore --define gotags=bazel \
		--test_env=BAZEL_WORKSPACE=$(shell bazel info workspace) \
		--sandbox_writable_path=$(shell bazel info workspace)/testdata \
		--sandbox_writable_path=$(shell bazel info workspace)/internal/testutils/parser/testdata \
		--sandbox_writable_path=$(shell bazel info workspace)/internal/testutils/parser/lexer/testdata" >> $@.tmp
	@echo "try-import %workspace%/.bazelrc.user" >> $@.tmp
	@mv $@.tmp .bazelrc
	@bazel run //:gazelle -- update-repos \
		-from_file=go.mod -prune=true \
		-build_file_proto_mode=disable_global \
		-to_macro=DEPS.bzl%go_deps &> /dev/null
	@bazel run //:gazelle &> /dev/null

gogenerate:
	@echo "--- generating go:generate files"
	@go generate ./...
	@echo "--- tidying go.{mod,sum}"
	@go mod tidy

generate: gogenerate swig proto bzl
	@echo "ok"

FORCE: ;
