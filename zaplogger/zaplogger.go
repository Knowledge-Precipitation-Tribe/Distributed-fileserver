package zaplogger

import "go.uber.org/zap"

var logger *zap.Logger

func init(){
	var err error
	logger, err = zap.NewProduction()
	if err != nil{
		panic(err)
	}
}

func GetLogger() *zap.Logger{
	if logger != nil{
		return logger
	}
	logger, err := zap.NewProduction()
	if err != nil{
		panic(err)
	}
	return logger
}

//输出日志到文件
func GetLoggerToFile(logfile string) *zap.Logger{
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{
		logfile,
	}
	loggerFile, err := cfg.Build()
	if err != nil{
		panic(err)
	}
	return loggerFile
}