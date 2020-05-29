package network

import (
	"net/http"
	"github.com/satori/go.uuid"
	"time"
	"fmt"
	"store-management-be/baselib/logger"
	"strconv"
	"strings"
	"store-management-be/database/redis"
)

type Handler func(requester *Requester, w http.ResponseWriter, r *http.Request)

func RecoveryHandler(handler Handler) Handler {
	logger.LogSugar.Infof("-----------RecoveryHandler-----------")
	return func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		defer func() {
			if err := recover(); err != nil {
				errMsg := fmt.Sprintf("Panic: %v", err)
				logger.LogSugar.Errorf("unknown err:%s", errMsg)
				NewHttpFailure(http.StatusInternalServerError, unknownPanic).AppendErr(errMsg).Response(w)
			}
		}()
		handler(requester, w, r)
	}
}

func RequesterHandler(handler Handler) http.Handler {
	logger.LogSugar.Infof("-----------RequesterHandler-----------")
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u := uuid.NewV4()
		uuid := u.String()
		aRequester := &Requester{
			UUID: uuid,
			//User:      UserInfo{},
			timestamp: time.Now(),
		}
		recorder := NewResponseRecorder(w)
		recorder.uuid = uuid
		handler(aRequester, recorder, r)
	})
}

func WebLoginTokenHandler(handler Handler) Handler {
	logger.LogSugar.Infof("-----------WebLoginTokenHandler-----------")
	return func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		t := requester.Token
		accToken, ok := t.(*AccessToken)
		if !ok {
			name := r.PostFormValue("name")
			pwd := r.PostFormValue("password")
			u := UserInfo(name)
			fmt.Println("u=", u)
			if u["id"] == "" || u["store_id"] == "" {
				logger.LogSugar.Errorf("WebLoginTokenHandler, user not find, name:%s", name)
				NewFailure(invalidAccount).AppendErr(" name:" + name).Response(w)
				return
			}
			id, _ := strconv.Atoi(u["id"])
			bId, _ := strconv.Atoi(u["business_id"])
			sId, _ := strconv.Atoi(u["store_id"])
			requester.Token = &AccessToken{
				ID:       id,
				Name:     name,
				Phone:    u["phone"],
				Password: pwd,
				BusiId:   bId,
				BusiType: u["busi_type"],
				StoreId:  sId,
			}
		} else {
			requester.Token = accToken
		}
		handler(requester, w, r)
	}
}

func UserInfo(name string) map[string]string {
	// todo 模拟从数据库取数据
	m := make(map[string]string)
	if name == "yangjun" {
		m["id"] = "1000001" // 用户id
		m["phone"] = "18356628848"
		m["business_id"] = "101"   // 商户id
		m["busi_type"] = "service" // 商户类型
		m["store_id"] = "20001"    // 门店id
	}
	return m
}

// NoJWTHandler is the middleware for initiating information in requester instance,
// when there is no json web token in request.
func NoJWTHandler(handler Handler) Handler {
	return Handler(func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		// 触点编号校验
		var delegateCode string
		if r.Method == http.MethodGet {
			delegateCode = r.URL.Query().Get("delegate_code")
		} else {
			delegateCode = r.PostFormValue("delegate_code")
		}
		if delegateCode == "" {
			logger.LogSugar.Errorf("JWTHandler, delegateCode err, delegateCode:%s", delegateCode)
			NewHttpFailure(http.StatusBadRequest, lackParameter).AppendErr("缺少参数: delegate_code").Response(w)
			return
		}

		delegate := SharedManager.GetDelegateByCode(delegateCode)
		if delegate.IsEqualTo(InvalidClient) {
			logger.LogSugar.Errorf("JWTHandler, delegateCode err, delegateCode:%s", delegateCode)
			NewHttpFailure(http.StatusBadRequest, invalidDelegate).AppendErr("delegate_code: " + delegateCode).Response(w)
			return
		}
		if recorder, ok := w.(*ResponseRecorder); ok {
			recorder.delegate = delegate
		}
		requester.From = delegate
		//requester.Token = &AccessToken{}
		handler(requester, w, r)
	})
}

// JWTHandler is the middleware for parsing json web token,
// which extract information from token, and inject the information into requester instance.
func JWTHandler(handler Handler) Handler {
	logger.LogSugar.Infof("-----------JWTHandler-----------")
	return Handler(func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		authorization := r.Header.Get("Authorization")
		if authorization == "" {
			logger.LogSugar.Errorf("JWTHandler, Authorization err, Authorization:%s", authorization)
			NewFailure(lackParameter).AppendErr("缺少参数: Authorization").Response(w)
			return
		}

		claims := ParseStandardClaims(authorization)
		fun := SharedManager.GetTokenGenerator(claims.Subject)
		if fun == nil {
			logger.LogSugar.Errorf("JWTHandler, Authorization err, Authorization:%s", authorization)
			NewFailure(lackParameter).AppendErr("token type not find:" + claims.Subject).Response(w)
			return
		}
		t := fun()
		err := t.ParseWith(authorization)
		if err != nil {
			logger.LogSugar.Errorf("JWTHandler, Authorization err, Authorization:%s", authorization)
			NewFailure(lackParameter).AppendErr("err:" + err.Error()).Response(w)
			return
		}
		requester.Token = t
		delegate := SharedManager.GetDelegateByCode(claims.Audience)
		if recorder, ok := w.(*ResponseRecorder); ok {
			recorder.delegate = delegate
		}
		requester.From = delegate

		handler(requester, w, r)
	})
}

func ForbidRequestTooOftenHandler(handler Handler) Handler {
	logger.LogSugar.Infof("-----------ForbidRequestTooOftenHandler-----------")
	return Handler(func(requester *Requester, w http.ResponseWriter, r *http.Request) {
		now := time.Now().UnixNano()
		diff := now
		defer func() {
			if diff <= 1*1e9 {
				logger.LogSugar.Errorf("ForbidRequestTooOftenHandler, diff:%d", diff)
				NewFailure(requestIsTooOften).Response(w)
				return
			} else {
				handler(requester, w, r)
			}
		}()

		addr := GetRemoteAddr(requester, r)
		uri := fmt.Sprintf("%v:%v", r.Method, r.URL.Path)
		id := "guest"
		if requester.Token != nil {
			id = fmt.Sprintf("%s", requester.Token.Identifier())
		}
		key := fmt.Sprintf("%v:%v:%v", id, addr, uri)
		lastTimestamp, err := GetLastTimestamp(key)
		if err != nil && err.Error() != "redigo: nil returned" {
			logger.LogSugar.Errorf("ForbidRequestTooOftenHandler, key:%s err:%s", key, err)
			return
		}
		diff = now - lastTimestamp
		SetLastTimestamp(key, now)
	})
}

func GetLastTimestamp(key string) (int64, error) {
	return redis.GetInstance().GetInt64(key)
}

func SetLastTimestamp(key string, value int64) error {
	_, err := redis.GetInstance().Set(key, value, 60)
	return err
}

// 获取远端请求IP
func GetRemoteAddr(requester *Requester, r *http.Request) string {
	if requester.From.IsNormalREST {
		return r.RemoteAddr
	} else {
		// 对支付宝小程序，无法获取请求方IP
		strs := strings.Split(r.Header.Get("X-Forwarded-For"), ",")
		if len(strs) >= 0 {
			return strs[0]
		} else {
			return r.RemoteAddr
		}
	}
}
