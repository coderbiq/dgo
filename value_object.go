package dgo

import "strconv"

// ValueObject model
type ValueObject interface {
	Equal(other ValueObject) bool
	String() string
}

// Identity model
type Identity interface {
	ValueObject
}

// StringID identity of string
type StringID string

// Equal assert equals
func (id StringID) Equal(other ValueObject) bool {
	otherStringID, ok := other.(StringID)
	if !ok {
		return false
	}
	return string(id) == string(otherStringID)
}

func (id StringID) String() string {
	return string(id)
}

// LongID identity of int64
type LongID int64

// Equal assert equals
func (id LongID) Equal(other ValueObject) bool {
	otherLongID, ok := other.(LongID)
	if !ok {
		return false
	}
	return int64(id) == int64(otherLongID)
}

func (id LongID) String() string {
	return strconv.FormatInt(int64(id), 10)
}
