package tsi1_test

import (
	"bytes"
	"reflect"
	"testing"

	"github.com/influxdata/influxdb/tsdb/engine/tsi1"
)

// Ensure measurement blocks can be written and opened.
func TestMeasurementBlockWriter(t *testing.T) {
	// Write 3 measurements to writer.
	mw := tsi1.NewMeasurementBlockWriter()
	mw.Add([]byte("foo"), 100, []uint32{1, 3, 4})
	mw.Add([]byte("bar"), 200, []uint32{2})
	mw.Add([]byte("baz"), 300, []uint32{5, 6})

	// Encode into buffer.
	var buf bytes.Buffer
	if n, err := mw.WriteTo(&buf); err != nil {
		t.Fatal(err)
	} else if n == 0 {
		t.Fatal("expected bytes written")
	}

	// Unmarshal into a TagSet.
	var blk tsi1.MeasurementBlock
	if err := blk.UnmarshalBinary(buf.Bytes()); err != nil {
		t.Fatal(err)
	}

	// Verify data in block.
	if e, ok := blk.Elem([]byte("foo")); !ok {
		t.Fatal("expected element")
	} else if e.Offset != 100 {
		t.Fatalf("unexpected offset: %v", e.Offset)
	} else if !reflect.DeepEqual(e.SeriesIDs(), []uint32{1, 3, 4}) {
		t.Fatalf("unexpected series data: %#v", e.SeriesIDs())
	}

	if e, ok := blk.Elem([]byte("bar")); !ok {
		t.Fatal("expected element")
	} else if e.Offset != 200 {
		t.Fatalf("unexpected offset: %v", e.Offset)
	} else if !reflect.DeepEqual(e.SeriesIDs(), []uint32{2}) {
		t.Fatalf("unexpected series data: %#v", e.SeriesIDs())
	}

	if e, ok := blk.Elem([]byte("baz")); !ok {
		t.Fatal("expected element")
	} else if e.Offset != 300 {
		t.Fatalf("unexpected offset: %v", e.Offset)
	} else if !reflect.DeepEqual(e.SeriesIDs(), []uint32{5, 6}) {
		t.Fatalf("unexpected series data: %#v", e.SeriesIDs())
	}

	// Verify non-existent measurement doesn't exist.
	if _, ok := blk.Elem([]byte("BAD_MEASUREMENT")); ok {
		t.Fatal("expected no element")
	}
}
