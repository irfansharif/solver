// Copyright 2010-2021 Google LLC
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

// This .i file exposes the CP-SAT model API. It was adapted from
// ortools/sat/{python,java,csharp}/sat.i.

%include "stdint.i"

%include "ortools/base/base.i"
%include "internal/proto.i"

%{
#include "ortools/sat/cp_model.pb.h"
#include "ortools/sat/sat_parameters.pb.h"
#include "ortools/sat/swig_helper.h"
#include "ortools/util/sorted_interval_list.h"
%}

%go_import("github.com/irfansharif/solver/internal/pb")

%module(directors="1") operations_research_sat

PROTO_INPUT(operations_research::sat::CpModelProto,
            pb.CpModelProto,
            model_proto);

PROTO_INPUT(operations_research::sat::SatParameters,
            pb.SatParameters,
            parameters);

PROTO_INPUT(operations_research::sat::IntegerVariableProto,
            pb.IntegerVariableProto,
            variable_proto);

PROTO_INPUT(operations_research::sat::CpSolverResponse,
            pb.CpSolverResponse,
            response);

PROTO2_RETURN(operations_research::sat::CpSolverResponse,
              pb.CpSolverResponse);

%ignoreall

%unignore operations_research;
%unignore operations_research::sat;

// Wrap the SolveWrapper class.
%unignore operations_research::sat::SolveWrapper;
%unignore operations_research::sat::SolveWrapper::AddLogCallback;        // unused
%unignore operations_research::sat::SolveWrapper::AddSolutionCallback;
%unignore operations_research::sat::SolveWrapper::ClearSolutionCallback; // unused
%unignore operations_research::sat::SolveWrapper::SetParameters;
%unignore operations_research::sat::SolveWrapper::Solve;
%unignore operations_research::sat::SolveWrapper::StopSearch;            // unused

// Wrap the CpSatHelper class.
%unignore operations_research::sat::CpSatHelper;
%unignore operations_research::sat::CpSatHelper::ModelStats;          // unused
%unignore operations_research::sat::CpSatHelper::SolverResponseStats; // unused
%unignore operations_research::sat::CpSatHelper::ValidateModel;
%unignore operations_research::sat::CpSatHelper::VariableDomain;      // unused

%feature("director") operations_research::sat::LogCallback; // unused
%unignore operations_research::sat::LogCallback;
%unignore operations_research::sat::LogCallback::~LogCallback;
%unignore operations_research::sat::LogCallback::NewMessage;

%feature("director") operations_research::sat::SolutionCallback;
%unignore operations_research::sat::SolutionCallback;
%unignore operations_research::sat::SolutionCallback::~SolutionCallback;
%unignore operations_research::sat::SolutionCallback::BestObjectiveBound;
%feature("nodirector") operations_research::sat::SolutionCallback::BestObjectiveBound;
%unignore operations_research::sat::SolutionCallback::HasResponse;
%feature("nodirector") operations_research::sat::SolutionCallback::HasResponse;
%unignore operations_research::sat::SolutionCallback::NumBinaryPropagations;
%feature("nodirector") operations_research::sat::SolutionCallback::NumBinaryPropagations;
%unignore operations_research::sat::SolutionCallback::NumBooleans;
%feature("nodirector") operations_research::sat::SolutionCallback::NumBooleans;
%unignore operations_research::sat::SolutionCallback::NumBranches;
%feature("nodirector") operations_research::sat::SolutionCallback::NumBranches;
%unignore operations_research::sat::SolutionCallback::NumConflicts;
%feature("nodirector") operations_research::sat::SolutionCallback::NumConflicts;
%unignore operations_research::sat::SolutionCallback::NumIntegerPropagations;
%feature("nodirector") operations_research::sat::SolutionCallback::NumIntegerPropagations;
%unignore operations_research::sat::SolutionCallback::ObjectiveValue;
%feature("nodirector") operations_research::sat::SolutionCallback::ObjectiveValue;
%unignore operations_research::sat::SolutionCallback::OnSolutionCallback;
%unignore operations_research::sat::SolutionCallback::Response;
%feature("nodirector") operations_research::sat::SolutionCallback::Response;
%unignore operations_research::sat::SolutionCallback::SolutionBooleanValue;
%feature("nodirector") operations_research::sat::SolutionCallback::SolutionBooleanValue;
%unignore operations_research::sat::SolutionCallback::SolutionIntegerValue;
%feature("nodirector") operations_research::sat::SolutionCallback::SolutionIntegerValue;
%unignore operations_research::sat::SolutionCallback::StopSearch;
%feature("nodirector") operations_research::sat::SolutionCallback::StopSearch;
%unignore operations_research::sat::SolutionCallback::UserTime;
%feature("nodirector") operations_research::sat::SolutionCallback::UserTime;
%unignore operations_research::sat::SolutionCallback::WallTime;
%feature("nodirector") operations_research::sat::SolutionCallback::WallTime;

%include "ortools/sat/swig_helper.h"

%unignoreall
