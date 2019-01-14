package points

import (
	"time"

	"github.com/coderbiq/dgo/model"
)

type (
	// Points 定义积分点数据模型
	Points uint

	// Account 定义积分账户外观
	Account interface {
		ID() model.LongID
		OwnerID() model.StringID
		Points() Points
		Deposit(points Points)
		Consume(points Points) error
	}

	// AccountReadModel 定义积分账户读模型
	AccountReadModel struct {
		ID              int64
		OwnerID         string
		Points          uint
		DepositedPoints uint
		ConsumedPoints  uint
		Logs            []AccountLog
	}

	// AccountLog 定义积分账户变更历史读模型
	AccountLog struct {
		ID      int64
		Acount  string
		Detail  map[string]interface{}
		Created time.Time
	}
)

type baseAccount struct {
	id      model.LongID
	ownerID model.StringID
	points  Points
}

func (a baseAccount) ID() model.LongID {
	return a.id
}

func (a baseAccount) OwnerID() model.StringID {
	return a.ownerID
}

func (a baseAccount) Points() Points {
	return a.points
}

// Inc 返回当前积分加上一个积分值后的新积分值
func (p Points) Inc(points Points) Points {
	return Points(uint(p) + uint(points))
}

// Dec 返回当前积分值减去一个积分值后的新积分值
func (p Points) Dec(points Points) Points {
	return Points(uint(p) - uint(points))
}

// GreaterThan 返回当前积分值是否比指定积分值大
func (p Points) GreaterThan(other Points) bool {
	return uint(p) > uint(other)
}

// Zero 返回当前积分值是否为零
func (p Points) Zero() bool {
	return p == 0
}
