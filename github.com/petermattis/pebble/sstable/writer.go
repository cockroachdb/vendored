// Copyright 2011 The LevelDB-Go and Pebble Authors. All rights reserved. Use
// of this source code is governed by a BSD-style license that can be found in
// the LICENSE file.

package sstable

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"math"

	"github.com/golang/snappy"
	"github.com/petermattis/pebble/internal/base"
	"github.com/petermattis/pebble/internal/crc"
	"github.com/petermattis/pebble/internal/rangedel"
)

// WriterMetadata holds info about a finished sstable.
type WriterMetadata struct {
	Size           uint64
	SmallestPoint  InternalKey
	SmallestRange  InternalKey
	LargestPoint   InternalKey
	LargestRange   InternalKey
	SmallestSeqNum uint64
	LargestSeqNum  uint64
}

func (m *WriterMetadata) updateSeqNum(seqNum uint64) {
	if m.SmallestSeqNum > seqNum {
		m.SmallestSeqNum = seqNum
	}
	if m.LargestSeqNum < seqNum {
		m.LargestSeqNum = seqNum
	}
}

func (m *WriterMetadata) updateLargestPoint(key InternalKey) {
	// Avoid the memory allocation in InternalKey.Clone() by reusing the buffer.
	m.LargestPoint.UserKey = append(m.LargestPoint.UserKey[:0], key.UserKey...)
	m.LargestPoint.Trailer = key.Trailer
}

// Smallest returns the smaller of SmallestPoint and SmallestRange.
func (m *WriterMetadata) Smallest(cmp Compare) InternalKey {
	if m.SmallestPoint.UserKey == nil {
		return m.SmallestRange
	}
	if m.SmallestRange.UserKey == nil {
		return m.SmallestPoint
	}
	if base.InternalCompare(cmp, m.SmallestPoint, m.SmallestRange) < 0 {
		return m.SmallestPoint
	}
	return m.SmallestRange
}

// Largest returns the larget of LargestPoint and LargestRange.
func (m *WriterMetadata) Largest(cmp Compare) InternalKey {
	if m.LargestPoint.UserKey == nil {
		return m.LargestRange
	}
	if m.LargestRange.UserKey == nil {
		return m.LargestPoint
	}
	if base.InternalCompare(cmp, m.LargestPoint, m.LargestRange) > 0 {
		return m.LargestPoint
	}
	return m.LargestRange
}

type flusher interface {
	Flush() error
}

type writeCloseSyncer interface {
	io.WriteCloser
	Sync() error
}

// Writer is a table writer.
type Writer struct {
	writer    io.Writer
	bufWriter *bufio.Writer
	syncer    writeCloseSyncer
	meta      WriterMetadata
	err       error
	// The following fields are copied from Options.
	blockSize          int
	blockSizeThreshold int
	compare            Compare
	split              Split
	compression        Compression
	separator          Separator
	successor          Successor
	tableFormat        TableFormat
	// With two level indexes, the index/filter of a SST file is partitioned into
	// smaller blocks with an additional top-level index on them. When reading an
	// index/filter, only the top-level index is loaded into memory. The two level
	// index/filter then uses the top-level index to load on demand into the block
	// cache the partitions that are required to perform the index/filter query.
	//
	// Two level indexes are enabled automatically when there is more than one
	// index block.
	//
	// This is useful when there are very large index blocks, which generally occurs
	// with the usage of large keys. With large index blocks, the index blocks fight
	// the data blocks for block cache space and the index blocks are likely to be
	// re-read many times from the disk. The top level index, which has a much
	// smaller memory footprint, can be used to prevent the entire index block from
	// being loaded into the block cache.
	twoLevelIndex      bool
	// Internal flag to allow creation of range-del-v1 format blocks. Only used
	// for testing. Note that v2 format blocks are backwards compatible with v1
	// format blocks.
	rangeDelV1Format bool
	// A table is a series of blocks and a block's index entry contains a
	// separator key between one block and the next. Thus, a finished block
	// cannot be written until the first key in the next block is seen.
	// pendingBH is the blockHandle of a finished block that is waiting for
	// the next call to Set. If the writer is not in this state, pendingBH
	// is zero.
	pendingBH      blockHandle
	block          blockWriter
	indexBlock     blockWriter
	rangeDelBlock  blockWriter
	props          Properties
	propCollectors []TablePropertyCollector
	// compressedBuf is the destination buffer for snappy compression. It is
	// re-used over the lifetime of the writer, avoiding the allocation of a
	// temporary buffer for each block.
	compressedBuf []byte
	// filter accumulates the filter block. If populated, the filter ingests
	// either the output of w.split (i.e. a prefix extractor) if w.split is not
	// nil, or the full keys otherwise.
	filter filterWriter
	// tmp is a scratch buffer, large enough to hold either footerLen bytes,
	// blockTrailerLen bytes, or (5 * binary.MaxVarintLen64) bytes.
	tmp [rocksDBFooterLen]byte

	topLevelIndexBlock blockWriter
	indexPartitions    []blockWriter
}

// Set sets the value for the given key. The sequence number is set to
// 0. Intended for use to externally construct an sstable before ingestion into
// a DB.
//
// TODO(peter): untested
func (w *Writer) Set(key, value []byte) error {
	if w.err != nil {
		return w.err
	}
	return w.addPoint(base.MakeInternalKey(key, 0, InternalKeyKindSet), value)
}

// Delete deletes the value for the given key. The sequence number is set to
// 0. Intended for use to externally construct an sstable before ingestion into
// a DB.
//
// TODO(peter): untested
func (w *Writer) Delete(key []byte) error {
	if w.err != nil {
		return w.err
	}
	return w.addPoint(base.MakeInternalKey(key, 0, InternalKeyKindDelete), nil)
}

// DeleteRange deletes all of the keys (and values) in the range [start,end)
// (inclusive on start, exclusive on end). The sequence number is set to
// 0. Intended for use to externally construct an sstable before ingestion into
// a DB.
//
// TODO(peter): untested
func (w *Writer) DeleteRange(start, end []byte) error {
	if w.err != nil {
		return w.err
	}
	return w.addTombstone(base.MakeInternalKey(start, 0, InternalKeyKindRangeDelete), end)
}

// Merge adds an action to the DB that merges the value at key with the new
// value. The details of the merge are dependent upon the configured merge
// operator. The sequence number is set to 0. Intended for use to externally
// construct an sstable before ingestion into a DB.
//
// TODO(peter): untested
func (w *Writer) Merge(key, value []byte) error {
	if w.err != nil {
		return w.err
	}
	return w.addPoint(base.MakeInternalKey(key, 0, InternalKeyKindMerge), value)
}

// Add adds a key/value pair to the table being written. For a given Writer,
// the keys passed to Add must be in increasing order. The exception to this
// rule is range deletion tombstones. Range deletion tombstones need to be
// added ordered by their start key, but they can be added out of order from
// point entries. Additionally, range deletion tombstones must be fragmented
// (i.e. by rangedel.Fragmenter).
func (w *Writer) Add(key InternalKey, value []byte) error {
	if w.err != nil {
		return w.err
	}

	if key.Kind() == InternalKeyKindRangeDelete {
		return w.addTombstone(key, value)
	}
	return w.addPoint(key, value)
}

func (w *Writer) addPoint(key InternalKey, value []byte) error {
	if base.InternalCompare(w.compare, w.meta.LargestPoint, key) >= 0 {
		w.err = fmt.Errorf("pebble: keys must be added in order: %s, %s", w.meta.LargestPoint, key)
		return w.err
	}

	if err := w.maybeFlush(key, value); err != nil {
		return err
	}

	for i := range w.propCollectors {
		if err := w.propCollectors[i].Add(key, value); err != nil {
			return err
		}
	}

	w.meta.updateSeqNum(key.SeqNum())
	w.meta.updateLargestPoint(key)

	w.maybeAddToFilter(key.UserKey)

	if w.props.NumEntries == 0 {
		w.meta.SmallestPoint = key.Clone()
	}
	w.props.NumEntries++
	switch key.Kind() {
	case InternalKeyKindDelete:
		w.props.NumDeletions++
	case InternalKeyKindMerge:
		w.props.NumMergeOperands++
	}
	w.props.RawKeySize += uint64(key.Size())
	w.props.RawValueSize += uint64(len(value))
	w.block.add(key, value)
	return nil
}

func (w *Writer) addTombstone(key InternalKey, value []byte) error {
	if !w.rangeDelV1Format && w.rangeDelBlock.nEntries > 0 {
		// Check that tombstones are being added in fragmented order. If the two
		// tombstones overlap, their start and end keys must be identical.
		prevKey := base.DecodeInternalKey(w.rangeDelBlock.curKey)
		switch c := w.compare(prevKey.UserKey, key.UserKey); {
		case c > 0:
			w.err = fmt.Errorf("pebble: keys must be added in order: %s, %s", prevKey, key)
			return w.err
		case c == 0:
			prevValue := w.rangeDelBlock.curValue
			if w.compare(prevValue, value) != 0 {
				w.err = fmt.Errorf("pebble: overlapping tombstones must be fragmented: %s vs %s",
					rangedel.Tombstone{Start: prevKey, End: prevValue},
					rangedel.Tombstone{Start: key, End: value})
				return w.err
			}
			if prevKey.SeqNum() <= key.SeqNum() {
				w.err = fmt.Errorf("pebble: keys must be added in order: %s, %s", prevKey, key)
				return w.err
			}
		default:
			prevValue := w.rangeDelBlock.curValue
			if w.compare(prevValue, key.UserKey) > 0 {
				w.err = fmt.Errorf("pebble: overlapping tombstones must be fragmented: %s vs %s",
					rangedel.Tombstone{Start: prevKey, End: prevValue},
					rangedel.Tombstone{Start: key, End: value})
				return w.err
			}
		}
	}

	for i := range w.propCollectors {
		if err := w.propCollectors[i].Add(key, value); err != nil {
			return err
		}
	}

	w.meta.updateSeqNum(key.SeqNum())

	if w.props.NumRangeDeletions == 0 {
		w.meta.SmallestRange = key.Clone()
		w.meta.LargestRange = base.MakeRangeDeleteSentinelKey(value).Clone()
	} else if w.rangeDelV1Format {
		if base.InternalCompare(w.compare, w.meta.SmallestRange, key) > 0 {
			w.meta.SmallestRange = key.Clone()
		}
		end := base.MakeRangeDeleteSentinelKey(value)
		if base.InternalCompare(w.compare, w.meta.LargestRange, end) < 0 {
			w.meta.LargestRange = end.Clone()
		}
	}
	w.props.NumRangeDeletions++
	w.rangeDelBlock.add(key, value)
	return nil
}

func (w *Writer) maybeAddToFilter(key []byte) {
	if w.filter != nil {
		if w.split != nil {
			prefix := key[:w.split(key)]
			w.filter.addKey(prefix)
		} else {
			w.filter.addKey(key)
		}
	}
}

func (w *Writer) maybeFlush(key InternalKey, value []byte) error {
	if size := w.block.estimatedSize(); size < w.blockSize {
		// The block is currently smaller than the target size.
		if size <= w.blockSizeThreshold {
			// The block is smaller than the threshold size at which we'll consider
			// flushing it.
			return nil
		}
		newSize := size + key.Size() + len(value)
		if w.block.nEntries%w.block.restartInterval == 0 {
			newSize += 4
		}
		newSize += 4                              // varint for shared prefix length
		newSize += uvarintLen(uint32(key.Size())) // varint for unshared key bytes
		newSize += uvarintLen(uint32(len(value))) // varint for value size
		if newSize <= w.blockSize {
			// The block plus the new entry is smaller than the target size.
			return nil
		}
	}

	bh, err := w.finishBlock(&w.block)
	if err != nil {
		w.err = err
		return w.err
	}
	w.pendingBH = bh
	w.flushPendingBH(key)
	return nil
}

// flushPendingBH adds any pending block handle to the index entries.
func (w *Writer) flushPendingBH(key InternalKey) {
	if w.pendingBH.length == 0 {
		// A valid blockHandle must be non-zero.
		// In particular, it must have a non-zero length.
		return
	}
	prevKey := base.DecodeInternalKey(w.block.curKey)
	var sep InternalKey
	if key.UserKey == nil && key.Trailer == 0 {
		sep = prevKey.Successor(w.compare, w.successor, nil)
	} else {
		sep = prevKey.Separator(w.compare, w.separator, nil, key)
	}
	n := encodeBlockHandle(w.tmp[:], w.pendingBH)

	if w.indexBlock.estimatedSize() >= w.blockSize*(len(w.indexPartitions)+1) {
		// Enable two level indexes if there is more than one index block.
		// TODO(ryan): Change this to `true` and uncomment when the reader
		// is implemented.
		w.twoLevelIndex = false
		//w.finishIndexBlock()
	}

	w.indexBlock.add(sep, w.tmp[:n])

	w.pendingBH = blockHandle{}
}

// finishBlock finishes the current block and returns its block handle, which is
// its offset and length in the table.
func (w *Writer) finishBlock(block *blockWriter) (blockHandle, error) {
	bh, err := w.writeRawBlock(block.finish(), w.compression)

	// Calculate filters.
	if w.filter != nil {
		w.filter.finishBlock(w.meta.Size)
	}

	// Reset the per-block state.
	block.reset()
	return bh, err
}

// finishIndexBlock finishes the current index block and adds it to the top
// level index block. This is only used when two level indexes are enabled.
func (w *Writer) finishIndexBlock() {
	w.indexPartitions = append(w.indexPartitions, w.indexBlock)
	w.indexBlock = blockWriter{
		restartInterval: 1,
	}
}

func (w *Writer) writeTwoLevelIndex() (blockHandle, error) {
	// Add the final unfinished index.
	w.finishIndexBlock()

	for _, b := range w.indexPartitions {
		sep := base.DecodeInternalKey(b.curKey)
		bh, _ := w.writeRawBlock(b.finish(), w.compression)

		if w.filter != nil {
			w.filter.finishBlock(w.meta.Size)
		}

		n := encodeBlockHandle(w.tmp[:], bh)
		w.topLevelIndexBlock.add(sep, w.tmp[:n])

		w.props.IndexSize += uint64(len(b.buf))
	}

	// NB: RocksDB includes the block trailer length in the index size
	// property, though it doesn't include the trailer in the top level
	// index size property.
	w.props.IndexPartitions = uint64(len(w.indexPartitions))
	w.props.TopLevelIndexSize = uint64(w.topLevelIndexBlock.estimatedSize())
	w.props.IndexSize += w.props.TopLevelIndexSize + blockTrailerLen

	return w.finishBlock(&w.topLevelIndexBlock)
}

func (w *Writer) writeRawBlock(b []byte, compression Compression) (blockHandle, error) {
	blockType := noCompressionBlockType
	if compression == SnappyCompression {
		// Compress the buffer, discarding the result if the improvement isn't at
		// least 12.5%.
		compressed := snappy.Encode(w.compressedBuf, b)
		w.compressedBuf = compressed[:cap(compressed)]
		if len(compressed) < len(b)-len(b)/8 {
			blockType = snappyCompressionBlockType
			b = compressed
		}
	}
	w.tmp[0] = blockType

	// Calculate the checksum.
	checksum := crc.New(b).Update(w.tmp[:1]).Value()
	binary.LittleEndian.PutUint32(w.tmp[1:5], checksum)
	bh := blockHandle{w.meta.Size, uint64(len(b))}

	// Write the bytes to the file.
	n, err := w.writer.Write(b)
	if err != nil {
		return blockHandle{}, err
	}
	w.meta.Size += uint64(n)
	n, err = w.writer.Write(w.tmp[:blockTrailerLen])
	if err != nil {
		return blockHandle{}, err
	}
	w.meta.Size += uint64(n)

	return bh, nil
}

// Close finishes writing the table and closes the underlying file that the
// table was written to.
func (w *Writer) Close() (err error) {
	defer func() {
		if w.syncer == nil {
			return
		}
		err1 := w.syncer.Close()
		if err == nil {
			err = err1
		}
		w.syncer = nil
	}()
	if w.err != nil {
		return w.err
	}

	// Finish the last data block, or force an empty data block if there
	// aren't any data blocks at all.
	w.flushPendingBH(InternalKey{})
	if w.block.nEntries > 0 || w.indexBlock.nEntries == 0 {
		bh, err := w.finishBlock(&w.block)
		if err != nil {
			w.err = err
			return w.err
		}
		w.pendingBH = bh
		w.flushPendingBH(InternalKey{})
	}
	w.props.DataSize = w.meta.Size
	w.props.NumDataBlocks = uint64(w.indexBlock.nEntries)

	// Write the filter block.
	var metaindex rawBlockWriter
	metaindex.restartInterval = 1
	if w.filter != nil {
		b, err := w.filter.finish()
		if err != nil {
			w.err = err
			return w.err
		}
		bh, err := w.writeRawBlock(b, NoCompression)
		if err != nil {
			w.err = err
			return w.err
		}
		n := encodeBlockHandle(w.tmp[:], bh)
		metaindex.add(InternalKey{UserKey: []byte(w.filter.metaName())}, w.tmp[:n])
		w.props.FilterPolicyName = w.filter.policyName()
		w.props.FilterSize = bh.length
	}

	var indexBH blockHandle
	if w.twoLevelIndex {
		w.props.IndexType = twoLevelIndex
		// Write the two level index block.
		indexBH, err = w.writeTwoLevelIndex()
		if err != nil {
			w.err = err
			return w.err
		}
	} else {
		w.props.IndexType = binarySearchIndex
		// NB: RocksDB includes the block trailer length in the index size
		// property, though it doesn't include the trailer in the filter size
		// property.
		w.props.IndexSize = uint64(w.indexBlock.estimatedSize()) + blockTrailerLen

		// Write the single level index block.
		indexBH, err = w.finishBlock(&w.indexBlock)
		if err != nil {
			w.err = err
			return w.err
		}
	}

	// Write the range-del block.
	if w.props.NumRangeDeletions > 0 {
		if !w.rangeDelV1Format {
			// Because the range tombstones are fragmented, the end key of the last
			// added range tombstone will be the largest range tombstone key. Note
			// that we need to make this into a range deletion sentinel because
			// sstable boundaries are inclusive while the end key of a range deletion
			// tombstone is exclusive.
			w.meta.LargestRange = base.MakeRangeDeleteSentinelKey(w.rangeDelBlock.curValue)
		}
		b := w.rangeDelBlock.finish()
		bh, err := w.writeRawBlock(b, w.compression)
		if err != nil {
			w.err = err
			return w.err
		}
		n := encodeBlockHandle(w.tmp[:], bh)
		// The v2 range-del block encoding is backwards compatible with the v1
		// encoding. We add meta-index entries for both the old name and the new
		// name so that old code can continue to find the range-del block and new
		// code knows that the range tombstones in the block are fragmented and
		// sorted.
		metaindex.add(InternalKey{UserKey: []byte(metaRangeDelName)}, w.tmp[:n])
		if !w.rangeDelV1Format {
			metaindex.add(InternalKey{UserKey: []byte(metaRangeDelV2Name)}, w.tmp[:n])
		}
	}

	{
		userProps := make(map[string]string)
		for i := range w.propCollectors {
			if err := w.propCollectors[i].Finish(userProps); err != nil {
				return err
			}
		}
		if len(userProps) > 0 {
			w.props.UserProperties = userProps
		}

		// Write the properties block.
		var raw rawBlockWriter
		// The restart interval is set to infinity because the properties block
		// is always read sequentially and cached in a heap located object. This
		// reduces table size without a significant impact on performance.
		raw.restartInterval = propertiesBlockRestartInterval
		w.props.CompressionOptions = rocksDBCompressionOptions
		w.props.save(&raw)
		bh, err := w.writeRawBlock(raw.finish(), NoCompression)
		if err != nil {
			w.err = err
			return w.err
		}
		n := encodeBlockHandle(w.tmp[:], bh)
		metaindex.add(InternalKey{UserKey: []byte(metaPropertiesName)}, w.tmp[:n])
	}

	// Write the metaindex block. It might be an empty block, if the filter
	// policy is nil.
	metaindexBH, err := w.finishBlock(&metaindex.blockWriter)
	if err != nil {
		w.err = err
		return w.err
	}

	// Write the table footer.
	footer := footer{
		format:      w.tableFormat,
		checksum:    checksumCRC32c,
		metaindexBH: metaindexBH,
		indexBH:     indexBH,
	}
	var n int
	if n, err = w.writer.Write(footer.encode(w.tmp[:])); err != nil {
		w.err = err
		return w.err
	}
	w.meta.Size += uint64(n)

	// Flush the buffer.
	if w.bufWriter != nil {
		if err := w.bufWriter.Flush(); err != nil {
			w.err = err
			return err
		}
	}

	if err := w.syncer.Sync(); err != nil {
		w.err = err
		return err
	}

	// Make any future calls to Set or Close return an error.
	w.err = errors.New("pebble: writer is closed")
	return nil
}

// EstimatedSize returns the estimated size of the sstable being written if a
// called to Finish() was made without adding additional keys.
func (w *Writer) EstimatedSize() uint64 {
	return w.meta.Size + uint64(w.block.estimatedSize()+w.indexBlock.estimatedSize())
}

// Metadata returns the metadata for the finished sstable. Only valid to call
// after the sstable has been finished.
func (w *Writer) Metadata() (*WriterMetadata, error) {
	if w.syncer != nil {
		return nil, errors.New("pebble: writer is not closed")
	}
	return &w.meta, nil
}

// NewWriter returns a new table writer for the file. Closing the writer will
// close the file.
func NewWriter(f writeCloseSyncer, o *Options, lo TableOptions) *Writer {
	o = o.EnsureDefaults()
	lo = *lo.EnsureDefaults()

	w := &Writer{
		syncer: f,
		meta: WriterMetadata{
			SmallestSeqNum: math.MaxUint64,
		},
		blockSize:          lo.BlockSize,
		blockSizeThreshold: (lo.BlockSize*lo.BlockSizeThreshold + 99) / 100,
		compare:            o.Comparer.Compare,
		split:              o.Comparer.Split,
		compression:        lo.Compression,
		separator:          o.Comparer.Separator,
		successor:          o.Comparer.Successor,
		tableFormat:        o.TableFormat,
		block: blockWriter{
			restartInterval: lo.BlockRestartInterval,
		},
		indexBlock: blockWriter{
			restartInterval: 1,
		},
		rangeDelBlock: blockWriter{
			restartInterval: 1,
		},
		topLevelIndexBlock: blockWriter{
			restartInterval: 1,
		},
	}
	if f == nil {
		w.err = errors.New("pebble: nil file")
		return w
	}

	w.props.PrefixExtractorName = "nullptr"
	if lo.FilterPolicy != nil {
		switch lo.FilterType {
		case TableFilter:
			w.filter = newTableFilterWriter(lo.FilterPolicy)
			if w.split != nil {
				w.props.PrefixExtractorName = o.Comparer.Name
				w.props.PrefixFiltering = true
			} else {
				w.props.WholeKeyFiltering = true
			}
		default:
			panic(fmt.Sprintf("unknown filter type: %v", lo.FilterType))
		}
	}

	w.props.ColumnFamilyID = math.MaxInt32
	w.props.ComparatorName = o.Comparer.Name
	w.props.CompressionName = lo.Compression.String()
	w.props.MergeOperatorName = o.Merger.Name
	w.props.PropertyCollectorNames = "[]"
	w.props.Version = 2 // TODO(peter): what is this?

	if len(o.TablePropertyCollectors) > 0 {
		w.propCollectors = make([]TablePropertyCollector, len(o.TablePropertyCollectors))
		var buf bytes.Buffer
		buf.WriteString("[")
		for i := range o.TablePropertyCollectors {
			w.propCollectors[i] = o.TablePropertyCollectors[i]()
			if i > 0 {
				buf.WriteString(",")
			}
			buf.WriteString(w.propCollectors[i].Name())
		}
		buf.WriteString("]")
		w.props.PropertyCollectorNames = buf.String()
	}

	// If f does not have a Flush method, do our own buffering.
	if _, ok := f.(flusher); ok {
		w.writer = f
	} else {
		w.bufWriter = bufio.NewWriter(f)
		w.writer = w.bufWriter
	}
	return w
}
