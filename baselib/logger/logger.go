package logger

import (
	"encoding/json"
	"flag"
	"os"
	"runtime/debug"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"

	"fmt"
	"net"
	"net/http"
	//"winkim/baselib/base"
)

var (
	// empty
	emptyBytes = []byte("{}")

	Logger         *zap.Logger
	LogSugar       *zap.SugaredLogger
	gmt8OffsetSecs int64
	aLevel         zap.AtomicLevel
)

var logFile = flag.String("log_file", "", "If non-empty, write log files in this directory")

func init() {
	// _, offset := time.Now().Local().Zone()
	gmt8OffsetSecs = int64(time.Hour * 8)
}
func GMT8TimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Add(time.Hour * 8).Format("01-02 15:04:05.000"))
}
func LocalTimeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("01-02 15:04:05.000"))
}

func GetLoggers() (*zap.Logger, *zap.SugaredLogger) {
	return Logger, LogSugar
}

// logpath 日志文件路径
// loglevel 日志级别
func InitLogger(loglevel string, levelPort int) (*zap.Logger, *zap.SugaredLogger) {
	// var level zapcore.Level
	aLevel = zap.NewAtomicLevel()
	//动态调整日志等级，用法如下：
	//“debug” “info” “warn” “error”
	//curl -XPUT --data '{"level":"debug"}' http://localhost:9090/handle/level
	//查询日志等级	curl http://localhost:9090/handle/level		输出：{"level":"info"}
	if levelPort != 0 && LocalPortAvailable(levelPort) {
		fmt.Println("zap log levelPort", levelPort)
		http.HandleFunc("/handle/level", aLevel.ServeHTTP)
		go func() {
			if err := http.ListenAndServe(fmt.Sprintf("127.0.0.1:%d", levelPort), nil); err != nil {
				panic(err)
			}
		}()
	}

	switch loglevel {
	case "debug":
		aLevel.SetLevel(zap.DebugLevel)
	case "info":
		aLevel.SetLevel(zap.InfoLevel)
	case "warn":
		aLevel.SetLevel(zap.WarnLevel)
	case "error":
		aLevel.SetLevel(zap.ErrorLevel)
	default:
		aLevel.SetLevel(zap.InfoLevel)
	}

	opts := []zap.Option{}
	opts = append(opts, zap.AddCaller())

	if len(*logFile) > 0 {
		//	写文件
		hook := lumberjack.Logger{
			Filename:   *logFile, // 日志文件路径
			MaxSize:    200,      // megabytes	默认100M
			MaxBackups: 4,        // 最多保留*个备份
			MaxAge:     7,        //days
			Compress:   true,     // 是否压缩 disabled by default
		}

		encoderConfig := zap.NewProductionEncoderConfig()
		encoderConfig.EncodeTime = GMT8TimeEncoder
		w := zapcore.AddSync(&hook)
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			w,
			aLevel,
		)
		Logger = zap.New(core, opts...)
	} else {
		//写stdError
		encoderConfig := zap.NewDevelopmentEncoderConfig()
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
		encoderConfig.EncodeTime = LocalTimeEncoder

		consoleDebugging := zapcore.Lock(os.Stdout)
		core := zapcore.NewCore(
			zapcore.NewConsoleEncoder(encoderConfig),
			consoleDebugging,
			aLevel,
		)
		opts = append(opts, zap.Development())
		Logger = zap.New(core, opts...)
	}

	LogSugar = Logger.Sugar()

	LogSugar.Warnf("SugarLogger init success: %v, logFile:%s", time.Now(), *logFile)
	Logger.Info("DefaultLogger init success")

	return Logger, LogSugar
}

func JsonDebugData(message interface{}) []byte {
	defer func() {
		if err := recover(); err != nil {
			LogSugar.Errorf("recover error JsonDebugData panic: %v", err)
			LogSugar.Errorf("recover error JsonDebugData - %v", string(debug.Stack()))
		}
	}()

	if data, err := json.Marshal(message); err == nil {
		return data
	}
	return emptyBytes
}

//此端口当前是否可用 只针对 127.0.0.1
func LocalPortAvailable(port int) bool {
	return IpPortAvailable("127.0.0.1", port)
}

func IpPortAvailable(ip string, port int) bool {
	if port >= 65535 || port <= 1024 {
		return false
	}
	laddr, err := net.ResolveTCPAddr("tcp4", fmt.Sprintf("%s:%d", ip, port))
	if err != nil {
		return false
	}

	listener, err := net.ListenTCP("tcp4", laddr)
	if err != nil {
		return false
	} else {
		listener.Close()
		return true
	}
}
