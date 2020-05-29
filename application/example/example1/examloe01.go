package example1

import (
	"net/http"
	"store-management-be/baselib/network"
	"store-management-be/baselib/logger"
	"store-management-be/application/example/model"
	"store-management-be/database/redis"
)

func GetMsg(requester *network.Requester, w http.ResponseWriter, r *http.Request) {
	logger.LogSugar.Infof("url:%v", r.URL.Path)

	_,err := model.QryMasterData()
	if err != nil {
		panic(err)
	}

	d, err := redis.GetInstance().GetString("a")
	if err != nil {
		panic(err)
	}
	network.NewSuccess(d).Response(w)
	//network.NewFailure(defaultFail).Response(w)
	//network.NewHttpFailure(http.StatusInternalServerError, defaultFail).Response(w)
}
