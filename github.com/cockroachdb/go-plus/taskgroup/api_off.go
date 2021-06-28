// Copyright 2021 The Cockroach Authors.
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

// +build !goplus

package taskgroup

import "runtime/metrics"

type internalTaskGroup = interface{}

func internalSupported() bool { return false }

func internalSetTaskGroup() T {
	return nil
}

func internalReadTaskGroupMetrics(taskGroup T, m []metrics.Sample) {
	for i := range m {
		// Initialize the Value to have KindBad (KindBad is the
		// default value 0 for the kind field).
		m[i].Value = metrics.Value{}
	}
}
