package network

func init() {
	d1 := Delegate{
		Code:         "10001",
		CodeType:     "web01",
		Desc:         "web端后台管理",
		IsNormalREST: true,
	}
	//Delegate
	SharedManager = &NetworkManager{
		delegatesMap: map[string]Delegate{
			d1.Code: d1,
		},
	}

	SharedManager.TokenFactory = make(map[string]func() Tokenable)
	// 注册不同类型的token
	SharedManager.RegisterToken(func() Tokenable {
		return new(AccessToken)
	})
}
