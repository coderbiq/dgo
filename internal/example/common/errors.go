package common

import "fmt"

// PointsInsufficientError 描述积分点数不足错误
type PointsInsufficientError struct {
	Current Points
	Expect  Points
}

func (err PointsInsufficientError) Error() string {
	return fmt.Sprintf("积分不足 %d 点, 当前为 %d 点。", err.Expect, err.Current)
}
