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