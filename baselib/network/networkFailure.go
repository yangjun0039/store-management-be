package network

import "fmt"

type networkFailure int

const DiscardMsg = "服务器内部错误"

const (
	jsonSerializationFailure networkFailure = iota + 1
	xmlSerializationFailure
	invalidJSONWebToken
	lackParameter
	responseTimeout
	invalidAccount
	invalidpubKey
	invalidDelegate
	requestIsTooOften
	verificationIsFailed
	unknownPanic
	invalidIP
	secretKeyFailure
	permissionNotAllow
	defaultfailure
	parameterFail
	rsaDecrpptFail
	nonceVerifyFail
)

func (this networkFailure) Code() string {
	return "network-" + fmt.Sprintf("%04d", this)
}

func (this networkFailure) ErrorMsg() string {
	switch this {
	case jsonSerializationFailure:
		return "json序列化失败,"
	case xmlSerializationFailure:
		return "xml序列化失败,"
	case invalidJSONWebToken:
		return "invalid jwt,"
	case lackParameter:
		return "authorization is lacked,"
	case responseTimeout:
		return "response timeout,"
	case requestIsTooOften:
		return "request is too often,"
	case unknownPanic:
		return "未知的panic,"
	case invalidAccount:
		return "invalid Account,"
	case invalidpubKey:
		return "invalid PubKey,"
	case invalidDelegate:
		return "invalid delegate,"
	case verificationIsFailed:
		return "signature verification is failed,"
	case invalidIP:
		return "invalid IP,"
	case secretKeyFailure:
		return "secret key failure,"
	case permissionNotAllow:
		return "permission not allow,"
	case parameterFail:
		return "parameter error,"
	case rsaDecrpptFail:
		return "rsa decrypt error,"
	case nonceVerifyFail:
		return "nonce verify error,"
	default:
		return DiscardMsg
	}
}

func (this networkFailure) DisplayedMsg() string {
	switch this {
	case invalidJSONWebToken:
		return "无效的令牌"
	case lackParameter:
		return "缺少授权信息"
	case responseTimeout:
		return "服务器响应超时"
	case requestIsTooOften:
		return "请求过于频繁"
	case invalidAccount:
		return "无效的账号"
	case invalidpubKey:
		return "公钥出错"
	case invalidDelegate:
		return "无效的触点"
	case verificationIsFailed:
		return "签名验证失败"
	case invalidIP:
		return "invalid IP"
	case secretKeyFailure:
		return "获取jwt签名密钥失败"
	case permissionNotAllow:
		return "权限不被允许"
	case parameterFail:
		return "参数错误"
	case rsaDecrpptFail:
		return "无效的rsa密文"
	case nonceVerifyFail:
		return "nonce验证失败"
	default:
		return DiscardMsg
	}
}
