package utils

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger Logger

type Logger struct {
	infoLogger  *zap.SugaredLogger
	errorLogger *zap.SugaredLogger
	fatalLogger *zap.SugaredLogger
}

func Init() {
	var err error
	logger.infoLogger, err = newLogger("./log/info.txt")
	if err != nil {
		panic(err)
	}

	//配置error logger
	logger.errorLogger, err = newLogger("./log/error.txt")
	if err != nil {
		panic(err)
	}

	//配置fatal logger
	logger.fatalLogger, err = newLogger("./log/fatal.txt")
	if err != nil {
		panic(err)
	}
}

func newLogger(fp string) (*zap.SugaredLogger, error) {
	file, err := os.OpenFile(fp, os.O_WRONLY|os.O_APPEND, 0666)
	writer := zapcore.AddSync(file) // 输出文件
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间戳的格式
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)
	log := zap.New(core, zap.AddCaller())
	return log.Sugar(), err
}

func InfoPrintln(str string, userId uint) {
	var printStr string
	if userId != 0 {
		printStr = fmt.Sprintf(str+"\tcause by userId=%d", userId)
	} else {
		printStr = fmt.Sprintf(str + "\tcause by not login user")
	}
	logger.infoLogger.Infof(printStr)
}

func ErrorPrintln(str string, userId uint, stackByte []byte) {
	printStr := str
	if userId != 0 {
		printStr = fmt.Sprintf(str+"\tcause by userId=%d", userId)
	} else {
		printStr = fmt.Sprintf(str)
	}
	printStr += "\n" + string(stackByte)
	logger.errorLogger.Errorf(printStr)
}

func FatalPrintln(str string, userId uint, stackByte []byte) {
	printStr := str
	if userId != 0 {
		printStr = fmt.Sprintf(str+"\tcause by userId=%d", userId)
	} else {
		printStr = fmt.Sprintf(str)
	}
	printStr += "\n" + string(stackByte)
	logger.errorLogger.Errorf(printStr)
}

func InitLog() {
	Init()
}
