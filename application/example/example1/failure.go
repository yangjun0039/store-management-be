package example1

import (
	"fmt"
)

type exampleFailure int

const (
	defaultFail   exampleFailure = iota + 1
	dataFail
	dbOperaFail
	parameterFail
)

func (this exampleFailure) Code() string {
	return "example-" + fmt.Sprintf("%04d", this)
}

func (this exampleFailure) ErrorMsg() string {
	switch this {
	case defaultFail:
		return "server error"
	case dataFail:
		return "data error"
	case dbOperaFail:
		return "db error"
	case parameterFail:
		return "param error"
	default:
		return "error"
	}
}

func (this exampleFailure) DisplayedMsg() string {
	switch this {
	case parameterFail:
		return "参数错误"
	default:
		return "服务器内部错误"
	}
}
