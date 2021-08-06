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

// This file was sourced from github.com/AirspaceTechnologies/or-tools.

%include "ortools/base/base.i"

%{
#include <vector>
#include "ortools/base/integral_types.h"
%}

%{
_goslice_ arrayIntToSlice(const int (&arr)[], size_t count) {
    _goslice_ slice;
    int *go_arr = (int*)malloc(sizeof(int[count]));
    slice.array = go_arr;
    slice.len = slice.cap = count;

    for (int i = 0; i < count; i++) {
        go_arr[i] = arr[i];
    }

    return slice;
}
%}

%insert(go_header) %{
type swig_goslice struct { arr uintptr; n int; c int }
func swigCopyIntSlice(intSlice *[]int) []int {
    newSlice := make([]int, len(*intSlice))
    for i := range newSlice {
        newSlice[i] = (*intSlice)[i]
    }
    p := *(*swig_goslice)(unsafe.Pointer(intSlice))
    Swig_free(p.arr)
    return newSlice
}
%}

%define _VECTOR_AS_GO_SLICE(ns, name, goname, ref, deref)
%{
std::vector< ns name ref > name##SliceToVector (_goslice_ slice) {
    std::vector< ns name ref > v;
    for (int i = 0; i < slice.len; i++) {
        ns name ref a = (( ns name * ref )slice.array)[i];
        v.push_back(a);
    }
    return v;
}

_goslice_ vectorTo##name##Slice (const std::vector< ns name ref >& arr) {
    _goslice_ slice;
    size_t count = arr.size();
    ns name * ref go_arr = ( ns name * ref )malloc(sizeof( ns name ref) * count);
    slice.array = go_arr;
    slice.len = slice.cap = count;

    for (int i = 0; i < count; i++) {
        go_arr[i] = arr[i];
    }

    return slice;
}
%}

%insert(go_header) %{
func swigCopy##name##Slice (s *[] goname ) [] goname {
    newSlice := make([] goname, len(*s))
    for i := range newSlice {
        newSlice[i] = (*s)[i]
    }
    p := *(*swig_goslice)(unsafe.Pointer(s))
    Swig_free(p.arr)
    return newSlice
}
%}

%typemap(gotype) std::vector< ns name ref > "[] goname"

%typemap(in) std::vector< ns name ref > %{
    $1 = name##SliceToVector($input);
%}

%typemap(out) std::vector< ns name ref > %{
    $result = vectorTo##name##Slice ($1);
%}

%typemap(goout) std::vector< ns name ref > %{
    $result = swigCopy##name##Slice(&$1)
%}


%typemap(gotype) const std::vector< ns name ref >& "[] goname"
%typemap(imtype) const std::vector< ns name ref >& "[] goname"

%typemap(goin) const std::vector< ns name ref > & %{
    $result = $1
%}

%typemap(goargout) const std::vector< ns name ref > & %{
%}

%typemap(argout) const std::vector< ns name ref > & %{
%}

%typemap(in) const std::vector< ns name ref > & %{
    $*1_ltype $1_arr;
    $1_arr = name##SliceToVector ($input);
    $1 = &$1_arr;
%}

%typemap(out) const std::vector< ns name ref > & %{
    $result = vectorTo##name##Slice (*$1);
%}

%typemap(goout) const std::vector< ns name ref > & %{
    $result = swigCopy##name##Slice(&$1)
%}

%insert(go_header) %{
type name##SliceWithPointer struct {
    slice [] goname
    ptr uintptr
}
%}

%typemap(gotype) std::vector< ns name ref > * "*[] goname"
%typemap(imtype) std::vector< ns name ref > * "* name##SliceWithPointer"

%typemap(goin) std::vector< ns name ref > * %{
    var $1_var name##SliceWithPointer
    $result = &$1_var
    $result.slice = *$1
%}

%typemap(goargout) std::vector< ns name ref > * %{
    *$1 = swigCopyIntSlice((*[]int)(unsafe.Pointer($1_var.ptr)))
    Swig_free($1_var.ptr)
%}

%typemap(in) std::vector< ns name ref > * %{
    sliceWithPointer* $1_ptr = (sliceWithPointer*)$input;

    $*1_ltype $1_arr;
    $1_arr = name##SliceToVector($1_ptr->slice);
    $1 = &$1_arr;
%}

%typemap(argout) std::vector< ns name ref > * %{
    $1_ptr->ptr = (_goslice_*)malloc(sizeof(_goslice_));
    *$1_ptr->ptr = arrayIntToSlice(*$1, sizeof($1.slice)/sizeof(void*));
%}

%enddef

#define nothing
#define VECTOR_AS_GO_SLICE(name, goname) _VECTOR_AS_GO_SLICE(nothing, name, goname, nothing, *)
#define VECTOR_AS_GO_SLICE_NAMESPACE(ns, name, goname) _VECTOR_AS_GO_SLICE(ns, name, goname, nothing, *)
#define VECTOR_AS_GO_POINTER_SLICE_NAMESPACE(ns, name, goname) _VECTOR_AS_GO_SLICE(ns, name, goname, *, nothing)

VECTOR_AS_GO_SLICE(int, int)
VECTOR_AS_GO_SLICE(int64_t, int64)
