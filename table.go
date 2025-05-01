package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash/fnv"
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

// Add item to table. 'value' is assumed to be a valid JSON string
func (t *Table) Add(key, value string) {
	idx := calculateIndex(key, t.Meta.Cap)

	WriteMessage(value, t.Data, idx)

	t.Meta.Keys = append(t.Meta.Keys, key)
	t.Meta.Size = t.Meta.Size + 1

	t.commitMetadata()
}

// Read and return JSON value stored under key 'key'
func (t *Table) Read(key string) (string, error) {
	idx := calculateIndex(key, t.Meta.Cap)

	firstChar, err := ReadMessage(t.Data[idx:idx+bitsPerByte], 1)

	if firstChar != fmt.Sprintf("%c", openCurlyBrace) {
		return "", fmt.Errorf("no data for key %s", key)
	}

	if err != nil {
		return "", err
	}

	m := make(map[string]string)
	for i := idx + bitsPerByte; i < idx+bucketSizeBytes; i += bitsPerByte {

		char, err := ReadMessage(t.Data[i:i+bitsPerByte], 1)
		if err != nil {
			return "", err
		}

		if char == fmt.Sprintf("%c", closeCurlyBrace) {
			jsonString, err := ReadMessage(
				t.Data[idx:idx+i+bitsPerByte+bitsPerByte+bitsPerByte],
				((i+bitsPerByte)-idx)/bitsPerByte,
			)
			if err != nil {
				return "", err
			}

			err = json.Unmarshal([]byte(jsonString), &m)
			if err == nil {
				return jsonString, nil
			}
		}
	}

	return "", fmt.Errorf("no data or invalid data for key %s", key)
}

// Calculate an index in Table.Data to store value at based on hash of 'key'
func calculateIndex(key string, capacity int) int {
	hash := fnv.New64a()
	hash.Write([]byte(key))
	idx := int(hash.Sum64() % uint64(capacity))

	// Not really needed but we like round numbers :D
	for idx%8 != 0 {
		idx += 1
	}

	fmt.Printf("Index: %d\n", idx)

	return idx
}

// Write struct Meta updates to data
func (t *Table) commitMetadata() {
	md, err := json.Marshal(t.Meta)
	if err != nil {
		panic(err)
	}
	t.Data = WriteMessage(string(md), t.Data, 0)
}
