package network


type NetworkManager struct {
	// 跨域白名单
	whiteOrigins map[string]bool
	// 触点模组
	delegatesMap map[string]Delegate
	// Token生成工厂
	TokenFactory map[string]func() Tokenable
}

var SharedManager *NetworkManager

func (m *NetworkManager) GetDelegateByCode(code string) Delegate {
	delegate := m.delegatesMap[code]
	if delegate.IsEqualTo(InvalidClient) {
		return InvalidClient
	}
	return delegate
}

func (m *NetworkManager) RegisterToken(tokenGenerator func() Tokenable) {
	m.TokenFactory[tokenGenerator().Subject()] = tokenGenerator
}

func (m *NetworkManager) GetTokenGenerator(subject string) func() Tokenable {
	return m.TokenFactory[subject]
}
