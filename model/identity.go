package model

import (
	"strconv"

	"github.com/bwmarrin/snowflake"
)

// IdentityGenerator 存储当前使用的标识生成器
// 主要用于为消息生成唯一标识，也可用于为聚合或实体生成标识。
var IdentityGenerator identityGenerator = defIdentityGenerator

type identityGenerator func() Identity

// IDFromInterface 根据一个 Interface 生成一个 Identity 实例
func IDFromInterface(v interface{}) Identity {
	switch id := v.(type) {
	case string:
		return StringID(id)
	case int64:
		return LongID(id)
	}
	return nil
}

// StringID 字符串类型的标识模型
type StringID string

// Equal 判断提供的另一个标识是否与当前标识为同一个标识
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

// Empty 返回当前标识是否为一个空值
func (id StringID) Empty() bool {
	return string(id) == ""
}

// LongID 长整数类型的标识模型
type LongID int64

// Equal 判断提供的另一个标识是否与当前标识为同一个标识
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

// Empty 返回当前标识是否为一个空值
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
