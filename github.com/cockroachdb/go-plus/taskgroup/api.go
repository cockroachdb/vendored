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

// Package taskgroup introduces the “task group” abstraction
// in the Go runtime, as discussed here:
// https://github.com/cockroachdb/cockroach/pull/60589
//
// In short, task groups are groups of related goroutines
// that share some accounting properties inside the Go runtime.
// Every new goroutine spawned inside a task group is automatically
// captured by the task group (i.e. new goroutine inherit
// their parent's task group by default).
//
// This similar to, but orthogonal to, the Inheritable goroutine ID
// (IGID) provided by package 'proc'. Task groups refer to accounting
// state maintained inside the Go runtime; IGIDs enable client apps to
// attach their own state to goroutine groups. Also, client apps
// control the specific value of IGIDs, whereas they cannot control the
// specific address or identification of task groups - these are chosen
// by the runtime when calling SetInternalTaskGroup().
//
// Programs should use the Supported() API function to determine whether
// the Go runtime extension is available.
package taskgroup

import "runtime/metrics"

// Supported returns true iff the task group extension is available
// in the Go runtime.
func Supported() bool {
	return internalSupported()
}

// T represents a task group inside the Go runtime.
//
// T values are reference-like, are comparable and can be nil, however
// no guarantee is provided about their concrete type.
type T = internalTaskGroup

// SetTaskGroup creates a new task group and attaches it to the
// current goroutine. It is inherited by future children goroutines.
// Top-level goroutines that have not been set a task group
// share a global (default) task group.
//
// If the extension is not supported, the operation is a no-op
// and returns a nil value.
func SetTaskGroup() T {
	return internalSetTaskGroup()
}

// ReadTaskGroupMetrics reads the runtime metrics for the specified
// task group. This is similar to the Go standard metrics.Read() but
// reads metrics scoped to just one task group.
//
// The following metric names are supported:
//
//   /taskgroup/sched/ticks:ticks
//   /taskgroup/sched/cputime:nanoseconds
//   /taskgroup/heap/largeHeapUsage:bytes
//
// If the extension is not supported, the metric value will be
// initialized with metric.KindBad.
//
// The reason why tg metrics are served by a dedicated function is
// that this function only locks the data structures for the specified
// task group, not the entire runtime. This makes it generally cheaper
// to use in a concurrent environment than metrics.Read().
func ReadTaskGroupMetrics(taskGroup T, m []metrics.Sample) {
	internalReadTaskGroupMetrics(taskGroup, m)
}
