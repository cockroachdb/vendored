// Copyright 2018 The LevelDB-Go and Pebble Authors. All rights reserved. Use
// of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package pebble

import "github.com/petermattis/pebble/internal/base"

// TableInfo exports the base.TableInfo type.
type TableInfo = base.TableInfo

// CompactionInfo exports the base.CompactionInfo type.
type CompactionInfo = base.CompactionInfo

// FlushInfo exports the base.FlushInfo type.
type FlushInfo = base.FlushInfo

// ManifestCreateInfo exports the base.ManifestCreateInfo type.
type ManifestCreateInfo = base.ManifestCreateInfo

// ManifestDeleteInfo exports the base.ManifestDeleteInfo type.
type ManifestDeleteInfo = base.ManifestDeleteInfo

// TableDeleteInfo exports the base.TableDeleteInfo type.
type TableDeleteInfo = base.TableDeleteInfo

// TableIngestInfo exports the base.TableIngestInfo type.
type TableIngestInfo = base.TableIngestInfo

// WALCreateInfo exports the base.WALCreateInfo type.
type WALCreateInfo = base.WALCreateInfo

// WALDeleteInfo exports the base.WALDeleteInfo type.
type WALDeleteInfo = base.WALDeleteInfo

// EventListener exports the base.EventListener type.
type EventListener = base.EventListener

// MakeLoggingEventListener exports the base.MakeLoggingEventListener function.
func MakeLoggingEventListener(logger Logger) EventListener {
	return base.MakeLoggingEventListener(logger)
}
