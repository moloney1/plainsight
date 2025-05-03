package main

import (
	"fmt"
	"slices"
	"testing"
)

type MockHasher struct{}

func (m MockHasher) Write(p []byte) (int, error) { return len(p), nil }
func (m MockHasher) Sum(b []byte) []byte         { return b }
func (m MockHasher) Reset()                      {}
func (m MockHasher) Size() int                   { return 0 }
func (m MockHasher) BlockSize() int              { return 0 }
func (m MockHasher) Sum64() uint64               { return 552 }

func TestNewTablePositive(t *testing.T) {
	bytes := make([]byte, 1000)
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
	bytes := make([]byte, 1000)
	bytes, _ = WriteMessage(metadataJson, bytes, 0)

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
		{"TestTableFromBytesN1", notQuiteJson},
	}

	for _, tc := range tests {
		bytes := make([]byte, 1000)
		bytes, _ = WriteMessage(tc.dataToWrite, bytes, 0)
		_, err := TableFromBytes(bytes, MockHasher{})
		if err == nil {
			t.Errorf("expected error")
		}
		fmt.Println(err)
	}
}

func TestReadPositive(t *testing.T) {
	metadataJson := "{\"cap\":1234,\"size\":1,\"keys\":[\"someKey\"]}"
	entryJson := "{\"user\":\"yourName\",\"pass\":\"hunter2\"}"
	bytes := make([]byte, 2000)
	bytes, _ = WriteMessage(metadataJson, bytes, 0)
	bytes, _ = WriteMessage(entryJson, bytes, 552) // to match what the MockHasher returns

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
		bytes := make([]byte, 2000)
		bytes, _ = WriteMessage(tc.metadataToWrite, bytes, 0)
		bytes, _ = WriteMessage(tc.dataToWrite, bytes, 552) // to match what the MockHasher returns

		table, err := TableFromBytes(bytes, MockHasher{})
		if err != nil {
			t.Fatalf("failed to create table from test data: %v", err)
		}

		_, err = table.Read("someKey")
		if err == nil {
			t.Errorf("expected error")
		}
	}
}

func TestAddPositive(t *testing.T) {
	table, err := NewTable(make([]byte, 2000), MockHasher{})
	if err != nil {
		t.Fatalf("error creating new table: %v", err)
	}
	if err = table.Add("someKey", "someValue"); err != nil {
		t.Errorf("unexpected error: %v", err)
	}
	addedToMetaKeys := slices.Contains(table.Meta.Keys, "someKey")
	if addedToMetaKeys != true {
		t.Error("expected someKey to be added to table metadata.Keys")
	}

}
