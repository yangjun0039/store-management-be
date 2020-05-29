package validation

import (
	"net/http"
	"store-management-be/baselib/network"
)

func Login(requester *network.Requester, w http.ResponseWriter, r *http.Request) {


	// todo 验证密码

	token, err := requester.Token.SignedString()
	if err != nil {
		network.NewFailure(tokenFail).AppendErr(err.Error()).Response(w)
		return
	}
	respInfo := map[string]string{
		"token": token,
	}
	network.NewSuccess(respInfo).Response(w)
}

//func UserInfo(name string) map[string]string {
//	// todo 模拟从数据库取数据
//	m := make(map[string]string)
//	if name == "yangjun" {
//		m["id"] = "1000001" // 用户id
//		m["phone"] = "18356628848"
//		m["business_id"] = "101"   // 商户id
//		m["busi_type"] = "service" // 商户类型
//		m["store_id"] = "20001"    // 门店id
//	}
//	return m
//}

func raise(err error) {
	if err != nil {
		panic(err)
	}
}
