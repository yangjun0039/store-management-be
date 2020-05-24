package configer

import (
	"flag"
	"fmt"
	"github.com/BurntSushi/toml"
	"store-management-be/database/mysql"
	"store-management-be/database/redis"
)

var Conf Config

func init() {
	// 加载配置
	var confFile string
	flag.StringVar(&confFile, "c", "./configer/config.toml", " set api config file path")
	if _, err := toml.DecodeFile(confFile, &Conf); err != nil {
		panic(fmt.Errorf("Decode_Config Error: %s", err.Error()))
	}
}

type Config struct {
	Server struct {
		Port     int
		Protocol string
	}
	Key struct {
		Certificate string
		Private     string
	}

	Mysql []mysql.MySQLConfig
	Redis redis.RedisConfig
}
