package jrm1

import (
	"fmt"
)

const (
	ErrFDuplicateMetaDataField  = "duplicate meta data field: %v"
	ErrFMetaDataFieldIsNotFound = "meta data field is not found: %v"
)

// ResponseMetaData is a set of named fields containing meta information about
// an RPC function call. It can be useful for transmission of identifier of RPC
// request, time measurement and many other things not directly related to an
// RPC function call.
type ResponseMetaData map[string]any

// AddField adds a field to the set.
func (md *ResponseMetaData) AddField(key string, value any) (err error) {
	_, isDuplicate := (*md)[key]
	if isDuplicate {
		return fmt.Errorf(ErrFDuplicateMetaDataField, key)
	}

	(*md)[key] = value

	return nil
}

// AddFieldFast adds a field to the set and panics on error.
func (md *ResponseMetaData) AddFieldFast(key string, value any) {
	err := md.AddField(key, value)
	if err != nil {
		panic(err)
	}
}

// GetField reads a field of the set.
func (md *ResponseMetaData) GetField(key string) (value any) {
	return (*md)[key]
}

// RemoveField deletes a field from the set.
func (md *ResponseMetaData) RemoveField(key string) (err error) {
	_, isPresent := (*md)[key]
	if !isPresent {
		return fmt.Errorf(ErrFMetaDataFieldIsNotFound, key)
	}

	delete(*md, key)

	return nil
}
