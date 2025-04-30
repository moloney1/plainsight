package main

import (
	"encoding/json"
)

const metaSizeBytes = 512
const bucketSizeBytes = 1024

type metadata struct {
	Cap  int      `json:"cap"`
	Size int      `json:"size"`
	Keys []string `json:"keys"`
}

type Table struct {
	Meta metadata
	Data []uint8
}

// Return an new empty Table
func NewTable(bytes []uint8) *Table {

	m := metadata{
		Cap:  (len(bytes) - metaSizeBytes) / bucketSizeBytes,
		Size: 0,
		Keys: make([]string, 0),
	}

	t := Table{
		Meta: m,
		Data: bytes, // should it be bytes[m.Cap:] ?
	}

	t.commitMetadata()

	return &t
}

// Return a previously populated Table
func TableFromBytes(bytes []uint8) {}

func (t *Table) commitMetadata() {
	md, err := json.Marshal(t.Meta)
	if err != nil {
		panic(err)
	}
	// fmt.Println(string(md))
	t.Data = WriteMessage(string(md), t.Data)
}
