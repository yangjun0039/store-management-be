package network

import (
	"time"
	"github.com/dgrijalva/jwt-go"
	"fmt"
	"net/http"
)

type Requester struct {
	UUID      string
	User      UserLogin
	timestamp time.Time
	From      Delegate    // 来源触点，使用触点代理号(请求来源)
	Business  StoreSource // 商户信息
	Token     Tokenable   // 所用令牌
}

type UserLogin struct {
	Phone    string `json:"phone"`
	Name     string `json:"name"`
	Password string `json:"password"`
	//ID    string `json:"id"`
	//Phone string `json:"phone"`
}

type StoreSource struct {
	BusiId   int    `json:"business_id"`
	BusiType string `json:"busi_type"`
	StoreId  int    `json:"store_id"`
}

// 触点代理编号
type Delegate struct {
	Code     string
	CodeType string
	Desc     string
	IsNormalREST bool // 是否正常RESTful API模式
}

type Tokenable interface {
	SignedString() (string, error)
	ParseWith(signedString string) error
	VerifiableMsg(requester *Requester, r *http.Request) (string, error)
	Verify(msg string, sig string) error
	Subject() string
	Identifier() int
	Validation() string
	//PubKey() []byte
	BusinessCert() StoreSource
	UserCert() UserData
}

func (d Delegate) IsEqualTo(another Delegate) bool {
	return d.Code == another.Code
}

var InvalidClient = Delegate{"", "", "无效的来源", true}

func LastOfMonth(t time.Time) time.Time {
	currentYear, currentMonth, _ := t.Date()
	location := t.Location()
	firstOfMonth := time.Date(currentYear, currentMonth, 1, 23, 59, 59, 999999999, location)
	lastOfMonth := firstOfMonth.AddDate(0, 1, -1)
	return lastOfMonth
}

// 解析签发的token
func ParseJWTWith(signedString string, secretKey string, claims jwt.Claims) error {
	token, err := jwt.ParseWithClaims(signedString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return err
	}
	// token有效性验证
	if token.Valid {
		return nil
	} else {
		return fmt.Errorf("JWTSignedString[%v] fail to parse mobile claims with secret key[%v]", signedString, secretKey)
	}
}

func ParseStandardClaims(signedString string) *jwt.StandardClaims {
	claims := new(jwt.StandardClaims)
	jwt.ParseWithClaims(signedString, claims, func(token *jwt.Token) (interface{}, error) {
		return nil, nil
	})
	return claims
}
