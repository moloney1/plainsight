package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash"
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
	Hash hash.Hash64
}

// Return an new empty Table
func NewTable(bytes []uint8, hasher hash.Hash64) (*Table, error) {

	m := metadata{
		Cap:  (len(bytes) - metaSizeBytes) / bucketSizeBytes,
		Size: 0,
		Keys: make([]string, 0),
	}

	t := Table{
		Meta: m,
		Data: bytes,
		Hash: hasher,
	}

	if err := t.commitMetadata(); err != nil {
		return &Table{}, err
	}

	return &t, nil
}

// Return a previously populated Table
func TableFromBytes(bytes []uint8, hasher hash.Hash64) (*Table, error) {

	firstChar, err := ReadMessage(bytes[:bitsPerByte], 1)
	if err != nil {
		return &Table{}, err
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
				return &Table{}, err
			}

			err = json.Unmarshal([]byte(jsonString), &meta)
			if err == nil {
				return &Table{
					Meta: meta,
					Data: bytes,
					Hash: hasher,
				}, nil
			}
		}
	}
	return &Table{}, errors.New("data invalid")
}

// Add item to table. 'value' is assumed to be a valid JSON string
func (t *Table) Add(key, value string) error {
	idx := t.calculateIndex(key, t.Meta.Cap)

	var err error
	t.Data, err = WriteMessage(value, t.Data, idx)
	if err != nil {
		return err
	}

	t.Meta.Keys = append(t.Meta.Keys, key)
	t.Meta.Size = t.Meta.Size + 1

	t.commitMetadata()

	return nil
}

// Read and return JSON value stored under key 'key'
func (t *Table) Read(key string) (string, error) {
	idx := t.calculateIndex(key, t.Meta.Cap)

	firstChar, err := ReadMessage(t.Data[idx:idx+bitsPerByte], 1)
	if err != nil {
		return "", err
	}

	if firstChar != fmt.Sprintf("%c", openCurlyBrace) {
		return "", fmt.Errorf("no data for key %s", key)
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
func (t *Table) calculateIndex(key string, capacity int) int {
	t.Hash.Write([]byte(key))
	idx := int(t.Hash.Sum64() % uint64(capacity))
	t.Hash.Reset()

	// Not really needed but we like round numbers :D
	for idx%8 != 0 {
		idx += 1
	}

	fmt.Printf("Index: %d\n", idx)

	return idx
}

// Write struct Meta updates to data
func (t *Table) commitMetadata() error {
	md, err := json.Marshal(t.Meta)
	if err != nil {
		return err
	}
	t.Data, err = WriteMessage(string(md), t.Data, 0)
	if err != nil {
		return err
	}
	return nil
}
