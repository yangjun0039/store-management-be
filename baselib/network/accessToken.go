package network

import (
	"github.com/dgrijalva/jwt-go"
	"time"
	"net/http"
	"fmt"
)

// 后台登陆token
type AccessToken struct {
	ID           int    //识别ID，用户账号
	DelegateCode string // 触点编号

	Name      string // 用户名
	Phone     string //手机号
	Password  string // 用户密码
	PublicKey []byte //通信公钥
	//PriKey []byte //
	BusiId   int    // 商户id
	BusiType string // 商户类型
	StoreId  int    // 门店id
}

type UserData struct {
	Id       int    `json:"id"`
	Phone    string `json:"phone"`
	Password string `json:"password"`
	//PublicKey []byte `json:"public_key"`
}

type AccessTokenClaims struct {
	jwt.StandardClaims
	User     UserData    `json:"user"`
	Business StoreSource `json:"business"`
}

var SecretKey = []byte("1234567890qwertyuioplkjhgfdsazxcvbnm")

func (t *AccessToken) SignedString() (string, error) {
	//id := fmt.Sprintf("%v:%v:%v", t.Subject(), t.Identifier(), t.Validation())
	claims := AccessTokenClaims{
		jwt.StandardClaims{
			NotBefore: time.Now().Unix(),
			ExpiresAt: LastOfMonth(time.Now()).Unix(), // token月末过期
			//ExpiresAt:time.Now().Unix()+20,
			Audience: t.DelegateCode,
			Issuer:   "yangjun",
			Id:       fmt.Sprintf("%d", t.ID),
			Subject:  t.Subject(),
		},
		UserData{
			Id:       t.ID,
			Phone:    t.Phone,
			Password: t.Password,
		},
		StoreSource{
			BusiId:   t.BusiId,
			BusiType: t.BusiType,
			StoreId:  t.StoreId,
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	ss, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return ss, nil
}

func (t *AccessToken) ParseWith(signedString string) error {
	claims := new(AccessTokenClaims)
	err := ParseJWTWith(signedString, string(SecretKey), claims)
	if err != nil {
		return err
	}
	t.Phone = claims.User.Phone
	t.Password = claims.User.Password
	t.ID = claims.User.Id
	t.BusiId = claims.Business.BusiId
	t.BusiType = claims.Business.BusiType
	t.StoreId = claims.Business.StoreId
	t.DelegateCode = claims.Audience
	//t.PublicKey = claims.User.PublicKey
	return nil
}

func (t *AccessToken) VerifiableMsg(requester *Requester, r *http.Request) (string, error) {
	return "", nil
}

func (t *AccessToken) Verify(msg string, hexSig string) error {
	return nil
}

func (t *AccessToken) Subject() string {
	return "AccessToken"
}

func (t *AccessToken) Identifier() int {
	return t.ID
}

func (t *AccessToken) Validation() string {
	return t.Password
}

func (t *AccessToken) UserCert() UserData {
	return UserData{t.ID, t.Phone, t.Password}
}

func (t *AccessToken) BusinessCert() StoreSource {
	return StoreSource{t.BusiId, t.BusiType, t.StoreId,}
}
