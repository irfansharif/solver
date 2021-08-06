// Copyright 2010-2018 Google LLC
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// TODO(user): make this SWIG file comply with the SWIG style guide.

// This file was sourced from github.com/AirspaceTechnologies/or-tools.

%include "ortools/base/base.i"

%include "internal/vector.i"

%{
#include <vector>
#include "ortools/base/integral_types.h"
%}

%go_import("fmt")
%go_import("unsafe")
%go_import("github.com/golang/protobuf/proto")

// SWIG macros to be used in generating Go wrappers for C++ protocol
// message parameters.  Each protocol message is serialized into
// byte[] before passing into (or returning from) C++ code.

// If the C++ function expects an input protocol message, transferring
// ownership to the caller (in C++):
//   foo(const MyProto* message,...)
// Use PROTO_INPUT macro:
//   PROTO_INPUT(MyProto, Google.Proto.Protos.Test.MyProto, message)
//
// if the C++ function returns a protocol message:
//   MyProto* foo();
// Use PROTO2_RETURN macro:
//   PROTO2_RETURN(MyProto, Google.Proto.Protos.Test.MyProto, true)
//
// Replace true by false if the C++ function returns a pointer to a
// protocol message object whose ownership is not transferred to the
// (C++) caller.
//
// Passing each protocol message from Go to C++ by value. Each ProtocolMessage
// is serialized into []byte when it is passed from Go to C++, the C++ code
// deserializes into C++ native protocol message.
//
// @param CppProtoType the fully qualified C++ protocol message type
// @param GoProtoType the corresponding fully qualified Go protocol message type
// @param param_name the parameter name
%define PROTO_INPUT(CppProtoType, GoProtoType, param_name)
%typemap(imtype) PROTO_TYPE* INPUT, PROTO_TYPE& INPUT "[]byte"
%typemap(gotype) PROTO_TYPE* INPUT, PROTO_TYPE& INPUT "GoProtoType"
%typemap(goin)   PROTO_TYPE* INPUT, PROTO_TYPE& INPUT {
  // go
  bytes, err := proto.Marshal(&$input)
  if err != nil {
    panic(fmt.Sprintf("Unable to convert input to []byte: %v", err))
  }
  $result = bytes
}
%typemap(in)     PROTO_TYPE* INPUT (CppProtoType temp), PROTO_TYPE& INPUT (CppProtoType temp) {
  // c
  bool parsed_ok = temp.ParseFromArray($input.array, $input.len);
  if (!parsed_ok) {
    _swig_gopanic("Unable to parse CppProtoType protocol message.");
  }
  $1 = &temp;
}
%apply PROTO_TYPE& INPUT { const CppProtoType& param_name }
%apply PROTO_TYPE& INPUT { CppProtoType& param_name }
%apply PROTO_TYPE* INPUT { const CppProtoType* param_name }
%apply PROTO_TYPE* INPUT { CppProtoType* param_name }
%enddef // end PROTO_INPUT

%define PROTO2_RETURN(CppProtoType, GoProtoType)
%typemap(imtype) CppProtoType "[]byte"
%typemap(gotype) CppProtoType "GoProtoType"
%typemap(goout)  CppProtoType {
  // go
  if err := proto.Unmarshal($1, &$result); err != nil {
    panic(fmt.Sprintf("Unable to parse GoProtoType protocol message: %v", err))
  }

  // free dynamic mem
  p := *(*swig_goslice)(unsafe.Pointer(&$1))
  Swig_free(p.arr)
}
%typemap(out) CppProtoType {
  uint8_t *go_arr = (uint8_t*)malloc($1.ByteSizeLong());
  $1.SerializeToArray(go_arr, $1.ByteSizeLong());

  _goslice_ slice;
  slice.array = go_arr;
  slice.len = slice.cap = $1.ByteSizeLong();
  $result = slice;
}
%enddef // end PROTO2_RETURN

