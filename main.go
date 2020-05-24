package main

import (
	"store-management-be/baselib/network"
	"store-management-be/application/example"
	"store-management-be/configer"
	"go.uber.org/zap"
	"store-management-be/baselib/logger"
	"store-management-be/database/mysql"
	"store-management-be/database/redis"
)

var (
	Log   *zap.Logger
	Sugar *zap.SugaredLogger
)

func main() {
	//日志等级和可通过XCURL 修改日志等级的端口
	Log, Sugar = logger.InitLogger("debug", 9087)

	initConfig()

	var httpProt = network.HTTPS
	protocol := configer.Conf.Server.Protocol
	if protocol == "HTTP" {
		httpProt = network.HTTP
	}

	start(httpProt, configer.Conf.Server.Port)
}

func initConfig() {
	// init mysql
	mysql.InitMysql(configer.Conf.Mysql)
	// init redis
	redis.InitRedis(configer.Conf.Redis)
}

func start(protocol network.NetProtocol, port int) {
	Sugar.Info("start server")
	r := network.NewRouter(nil)
	example.RouterRegisterMethods(r)

	//http.ListenAndServe(addr, r)
	r.Startup(protocol, uint64(port))
}
