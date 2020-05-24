package network

import (
	"net/http"
	"github.com/satori/go.uuid"
	"time"
	"fmt"
	"store-management-be/baselib/logger"
)

type Handler func(requester *Requester, w http.ResponseWriter, r *http.Request)

func RecoveryHandler(handler Handler) Handler {
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
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		u, _ := uuid.NewV4()
		uuid := u.String()
		aRequester := &Requester{
			UUID:      uuid,
			User:      UserInfo{},
			timestamp: time.Now(),
		}
		recorder := NewResponseRecorder(w)
		recorder.uuid = uuid
		handler(aRequester, recorder, r)
	})
}
