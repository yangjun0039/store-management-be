package member

import (
	"fmt"
)

type memberFailure int


const (
	defaultFail   memberFailure = iota + 1
	dataFail
	dbOperaFail
	parameterFail
	tokenFail

	NameNilFail
	PhoneFail
	MemLevelFail
	EffectTimeFail
)

func (this memberFailure) Code() string {
	return "member-" + fmt.Sprintf("%04d", this)
}

func (this memberFailure) ErrorMsg() string {
	switch this {
	case defaultFail:
		return "server error "
	case dataFail:
		return "data error "
	case dbOperaFail:
		return "db error "
	case parameterFail:
		return "param error "
	case tokenFail:
		return "get token error "
	case NameNilFail:
		return "name is nil error"
	case PhoneFail:
		return "phone invalid error"
	case MemLevelFail:
		return "member level error"
	case EffectTimeFail:
		return "effect time error"
	default:
		return "error "
	}
}

func (this memberFailure) DisplayedMsg() string {
	//if this.CustomMsg != "" {
	//	return this.CustomMsg
	//}
	switch this {
	case parameterFail:
		return "参数错误"
	case NameNilFail:
		return "用户名不能为空"
	case PhoneFail:
		return "手机号错误"
	case MemLevelFail:
		return "会员等级错误"
	case EffectTimeFail:
		return "会员有效时间错误"
	default:
		return "服务器内部错误"
	}
}

func (this custMemberFailure) CustomDisdMsg(msg string) {
	this.CustomMsg = msg
}

// custom
