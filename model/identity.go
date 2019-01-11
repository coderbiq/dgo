package model

import (
	"strconv"

	"github.com/bwmarrin/snowflake"
)

// IdentityGenerator generate identity
var IdentityGenerator identityGenerator = defIdentityGenerator

type identityGenerator func() Identity

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

// Empty assert empty
func (id StringID) Empty() bool {
	return string(id) == ""
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

// Empty assert empty
func (id LongID) Empty() bool {
	return int64(id) == 0
}

func defIdentityGenerator() Identity {
	node, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}
	return LongID(node.Generate())
}
