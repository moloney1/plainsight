package table

import (
	"encoding/json"
	"errors"
	"fmt"
	"hash"
	"math/rand"
	"slices"

	"github.com/moloney1/plainsight/internal/codec"
	"github.com/moloney1/plainsight/internal/constants"
)

const metaSizeBytes = 2048
const bucketSizeBytes = 4096

const openCurlyBrace = 123
const closeCurlyBrace = 125

const bitsPerByte = constants.BitsPerByte

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

type UserPass struct {
	User string `json:"user"`
	Pass string `json:"pass"`
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

	firstChar, err := codec.ReadMessage(bytes[:bitsPerByte], 1)
	if err != nil {
		return &Table{}, err
	}

	if firstChar != fmt.Sprintf("%c", openCurlyBrace) {
		return &Table{}, errors.New("no data found")
	}

	meta := metadata{}

	for i := bitsPerByte; i < metaSizeBytes; i += bitsPerByte { // TODO bounds?
		char, err := codec.ReadMessage(bytes[i:i+bitsPerByte], 1)
		if err != nil {
			return &Table{}, err
		}
		if char == fmt.Sprintf("%c", closeCurlyBrace) {

			jsonString, err := codec.ReadMessage(bytes[:i+bitsPerByte], (i+bitsPerByte)/bitsPerByte)
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

// Insert item to table. 'value' is assumed to be a valid JSON string
func (t *Table) Insert(key, value string) error {
	idx := t.calculateIndex(key, t.Meta.Cap)

	var err error
	t.Data, err = codec.WriteMessage(value, t.Data, idx)
	if err != nil {
		return err
	}

	t.Meta.Keys = append(t.Meta.Keys, key)
	t.Meta.Size = t.Meta.Size + 1

	if err := t.commitMetadata(); err != nil {
		return err
	}

	return nil
}

// Encode username/password pair to JSON and add to table
func (t *Table) AddUsernamePasswordPair(key, username, password string) error {

	pair, err := json.Marshal(UserPass{User: username, Pass: password})
	if err != nil {
		return err
	}

	return t.Insert(key, string(pair))
}

// Remove writes random data over where key is currently stored
func (t *Table) Remove(key string) error {

	if !slices.Contains(t.Meta.Keys, key) {
		return fmt.Errorf("nothing to delete at key %s", key)
	}

	idx := t.calculateIndex(key, t.Meta.Cap)

	// Remove the key from the Metadata
	t.Meta.Keys = slices.DeleteFunc(t.Meta.Keys, func(s string) bool {
		return s == key
	})
	t.Meta.Size -= 1
	if err := t.commitMetadata(); err != nil {
		return err
	}

	// Generate a random string that will take up exactly 'bucketSizeBytes'
	randBytes := make([]byte, bucketSizeBytes/bitsPerByte)
	for i := range len(randBytes) {
		randBytes[i] = byte(rand.Intn(128))
	}
	var err error
	t.Data, err = codec.WriteMessage(string(randBytes), t.Data, idx)
	if err != nil {
		return err
	}

	return nil
}

// Read and return JSON value stored under key 'key'
func (t *Table) Read(key string) (string, error) {
	idx := t.calculateIndex(key, t.Meta.Cap)

	firstChar, err := codec.ReadMessage(t.Data[idx:idx+bitsPerByte], 1)
	if err != nil {
		return "", err
	}

	if firstChar != fmt.Sprintf("%c", openCurlyBrace) {
		return "", fmt.Errorf("no data for key %s", key)
	}

	m := make(map[string]any)
	for i := idx + bitsPerByte; i < idx+bucketSizeBytes; i += bitsPerByte {

		char, err := codec.ReadMessage(t.Data[i:i+bitsPerByte], 1)
		if err != nil {
			return "", err
		}

		if char == fmt.Sprintf("%c", closeCurlyBrace) {
			jsonString, err := codec.ReadMessage(
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

// Return list of stored keys
func (t *Table) List() []string {
	return t.Meta.Keys
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

	return idx
}

// Write struct Meta updates to data
func (t *Table) commitMetadata() error {
	md, err := json.Marshal(t.Meta)
	if err != nil {
		return err
	}

	mdString := string(md)
	if len(mdString)*bitsPerByte > metaSizeBytes {
		return fmt.Errorf("cannot add: metadata section full, try removing some keys")
	}

	t.Data, err = codec.WriteMessage(string(md), t.Data, 0)
	if err != nil {
		return err
	}
	return nil
}
