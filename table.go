package main

import (
	"encoding/json"
	"errors"
	"fmt"
)

const metaSizeBytes = 512
const bucketSizeBytes = 1024

const openCurlyBrace = 123
const closeCurlyBrace = 125

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
		Data: bytes,
	}

	t.commitMetadata()

	return &t
}

// Return a previously populated Table
func TableFromBytes(bytes []uint8) (*Table, error) {

	firstChar, err := ReadMessage(bytes[:bitsPerByte], 1)
	if err != nil {
		panic(err)
	}

	if firstChar != fmt.Sprintf("%c", openCurlyBrace) {
		return &Table{}, errors.New("no data found")
	}

	meta := metadata{}

	for i := bitsPerByte; i < metaSizeBytes; i += bitsPerByte { // TODO find somewhere else for bitsPerByte const // TODO bounds?
		char, err := ReadMessage(bytes[i:i+bitsPerByte], 1)
		if err != nil {
			return &Table{}, err
		}
		if char == fmt.Sprintf("%c", closeCurlyBrace) {

			jsonString, err := ReadMessage(bytes[:i+bitsPerByte], (i+bitsPerByte)/bitsPerByte)
			if err != nil {
				panic(err)
			}

			err = json.Unmarshal([]byte(jsonString), &meta)
			if err == nil {
				return &Table{
					Meta: meta,
					Data: bytes,
				}, nil
			}
		}
	}
	return &Table{}, errors.New("data invalid")
}

func (t *Table) commitMetadata() {
	md, err := json.Marshal(t.Meta)
	if err != nil {
		panic(err)
	}
	t.Data = WriteMessage(string(md), t.Data, 0)
}
