// Copyright 2021 Irfan Sharif.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or
// implied. See the License for the specific language governing
// permissions and limitations under the License.

package solver

import (
	"fmt"
	"io"
	"log"
	"time"

	"github.com/irfansharif/solver/internal"
	"github.com/irfansharif/solver/internal/pb"
)

type Option func(o *options, s internal.SolveWrapper)

type options struct {
	params   pb.SatParameters
	logger   *log.Logger
	solution *solutionCallback
}

func (o *options) validate() (bool, error) {
	if o.params.GetEnumerateAllSolutions() && o.params.GetNumSearchWorkers() > 1 {
		return false, fmt.Errorf("cannot enumerate with parallelism > 1")
	}
	return true, nil
}

// WithTimeout configures the solver with a hard time limit.
func WithTimeout(d time.Duration) Option {
	return func(o *options, _ internal.SolveWrapper) {
		seconds := d.Seconds()
		o.params.MaxTimeInSeconds = &seconds
	}
}

// WithLogger configures the solver to route its internal logging to the given
// io.Writer, using the given prefix.
func WithLogger(w io.Writer, prefix string) Option {
	return func(o *options, s internal.SolveWrapper) {
		logSearchProgress, logToResponse, logToStdout := true, true, false
		o.params.LogSearchProgress = &logSearchProgress
		o.params.LogToStdout = &logToStdout
		o.params.LogToResponse = &logToResponse

		// TODO(irfansharif): Right now we're simply logging to the response
		// proto, which isn't being streamed during the search process and not
		// super. OR-Tools v9.0 does support an experimental logger callback
		// (looks identical to the solution callback), but that didn't work.
		//
		// Worth checking back on at some point.
		// https://github.com/google/or-tools/issues/1903
		o.logger = log.New(w, prefix, 0)
	}
}

// WithParallelism configures the solver to use the given number of parallel
// workers during search. If the number provided is <= 1, there will be no
// parallelism.
func WithParallelism(parallelism int) Option {
	return func(options *options, w internal.SolveWrapper) {
		threads := int32(parallelism)
		options.params.NumSearchWorkers = &threads
	}
}

// WithEnumeration configures the solver to enumerate over all solutions without
// objective. This option is incompatible with a parallelism greater than 1.
func WithEnumeration(f func(Result)) Option {
	return func(o *options, s internal.SolveWrapper) {
		enumerate := true
		o.params.EnumerateAllSolutions = &enumerate

		o.solution = &solutionCallback{f: f}
		o.solution.hook = internal.NewDirectorSolutionCallback(o.solution)
		s.AddSolutionCallback(o.solution.hook)
	}
}

// solutionCallback is used to hook into the underlying solver during its search
// process. It's invoked whenever a solution is found.
type solutionCallback struct {
	f    func(Result)
	hook internal.SolutionCallback
}

func (p *solutionCallback) OnSolutionCallback() {
	proto := p.hook.Response()
	p.f(Result{pb: &proto})
}
