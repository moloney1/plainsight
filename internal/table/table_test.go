// TODO convert identifiers of instances of struct Table to tbl
package table

import (
	"slices"
	"testing"

	"github.com/moloney1/plainsight/internal/codec"
)

const mockIndex = 552
const testSliceSize = metaSizeBytes + 2*bucketSizeBytes

type MockHasher struct{}

func (m MockHasher) Write(p []byte) (int, error) { return len(p), nil }
func (m MockHasher) Sum(b []byte) []byte         { return b }
func (m MockHasher) Reset()                      {}
func (m MockHasher) Size() int                   { return 0 }
func (m MockHasher) BlockSize() int              { return 0 }
func (m MockHasher) Sum64() uint64               { return mockIndex }

func TestNewTablePositive(t *testing.T) {
	bytes := make([]byte, testSliceSize)
	_, err := NewTable(bytes, MockHasher{})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestNewTableNegative(t *testing.T) {
	bytes := make([]byte, 16) // not enough space
	_, err := NewTable(bytes, MockHasher{})
	if err == nil {
		t.Errorf("expected error")
	}
}

func TestTableFromBytesPositive(t *testing.T) {
	metadataJson := "{\"cap\":1234,\"size\":1,\"keys\":[\"someKey\"]}"
	bytes := make([]byte, testSliceSize)
	bytes, _ = codec.WriteMessage(metadataJson, bytes, 0)

	_, err := TableFromBytes(bytes, MockHasher{})
	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
}

func TestTableFromBytesNegative(t *testing.T) {
	nonsense := "this is not json data"
	notQuiteJson := "{\"cap\":1234,\"size\":1,\"keys"

	tests := []struct {
		name        string
		dataToWrite string
	}{
		{"TestTableFromBytesN1", nonsense},
		{"TestTableFromBytesN2", notQuiteJson},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bytes := make([]byte, testSliceSize)
			bytes, _ = codec.WriteMessage(tc.dataToWrite, bytes, 0)
			_, err := TableFromBytes(bytes, MockHasher{})
			if err == nil {
				t.Errorf("expected error")
			}
			t.Log(err)
		})
	}
}

func TestReadPositive(t *testing.T) {
	metadataJson := "{\"cap\":1234,\"size\":1,\"keys\":[\"someKey\"]}"
	entryJson := "{\"user\":\"yourName\",\"pass\":\"hunter2\"}"
	bytes := make([]byte, testSliceSize)
	bytes, _ = codec.WriteMessage(metadataJson, bytes, 0)
	bytes, _ = codec.WriteMessage(entryJson, bytes, mockIndex)

	table, err := TableFromBytes(bytes, MockHasher{})
	if err != nil {
		t.Fatalf("failed to create table from test data: %v", err)
	}

	got, err := table.Read("myKey")
	if err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	if got != entryJson {
		t.Errorf("expected %s, got %s", entryJson, got)
	}
}

func TestReadNegative(t *testing.T) {
	tests := []struct {
		name            string
		metadataToWrite string
		dataToWrite     string
	}{
		{"TestReadN1", "{\"cap\":1234,\"size\":0,\"keys\":[]}", "nothing here"},        // no data
		{"TestReadN2", "{\"cap\":1234,\"size\":1,\"keys\":[\"someKey\"]}", "{\"us..."}, // data corrupted, or "{" is there by fluke
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			bytes := make([]byte, testSliceSize)
			bytes, _ = codec.WriteMessage(tc.metadataToWrite, bytes, 0)
			bytes, _ = codec.WriteMessage(tc.dataToWrite, bytes, mockIndex) // to match what the MockHasher returns

			table, err := TableFromBytes(bytes, MockHasher{})
			if err != nil {
				t.Fatalf("failed to create table from test data: %v", err)
			}

			_, err = table.Read("someKey")
			if err == nil {
				t.Errorf("expected error")
			}
		})
	}
}

func TestInsert(t *testing.T) {
	table, err := NewTable(make([]byte, testSliceSize), MockHasher{})
	if err != nil {
		t.Fatalf("error creating new table: %v", err)
	}
	if err = table.Insert("someKey", "someValue"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	addedToMetaKeys := slices.Contains(table.Meta.Keys, "someKey")
	if !addedToMetaKeys {
		t.Error("expected someKey to be added to table metadata.Keys")
	}
}

func TestInsertNegative(t *testing.T) {
	table, err := NewTable(make([]byte, testSliceSize), MockHasher{})
	if err != nil {
		t.Fatalf("error creating new table: %v", err)
	}
	longKey := "this-key-is-long-enough-to-overflow-the-metadata-section-this-key-is-long-enough-to-overflow-the-metadata-section-this-key-is-long-enough-to-overflow-the-metadata-section-this-key-is-long-enough-to-overflow-the-metadata-section-this-key-is-long-enough-to-overflow-the-metadata-section-this-key-is-long-enough-to-overflow-the-metadata-section-this-key-is-long-enough-to-overflow-the-metadata-section-this-key-is-long-enough-to-overflow-the-metadata-section-"
	if err = table.Insert(longKey, "someValue"); err == nil {
		t.Error("expected error when metadata section overflows")
	}
}

func TestRemovePositive(t *testing.T) {
	metadataJson := "{\"cap\":1234,\"size\":1,\"keys\":[\"someKey\"]}"
	entryJson := "{\"user\":\"yourName\",\"pass\":\"hunter2\"}"
	bytes := make([]byte, testSliceSize)
	bytes, _ = codec.WriteMessage(metadataJson, bytes, 0)
	bytes, _ = codec.WriteMessage(entryJson, bytes, mockIndex)

	tbl, _ := TableFromBytes(bytes, MockHasher{})

	sizeBefore := tbl.Meta.Size

	err := tbl.Remove("someKey")

	if err != nil {
		t.Errorf("unexpected error %v", err)
	}
	if slices.Contains(tbl.Meta.Keys, "someKey") {
		t.Errorf("expected key 'someKey' to have been deleted from table metadata: %v", tbl.Meta)
	}
	sizeAfter := tbl.Meta.Size
	if sizeAfter != sizeBefore-1 {
		t.Errorf("expected size to be reduced, size is %d", sizeAfter)
	}

}
