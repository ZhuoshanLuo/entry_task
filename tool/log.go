package tool

import (
	"fmt"
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var logger Logger

type Logger struct {
	InfoLogger  *zap.SugaredLogger
	ErrorLogger *zap.SugaredLogger
	FatalLogger *zap.SugaredLogger
}

func Init() {
	//配置info logger
	file, _ := os.Create("./log/info.txt")
	writer := zapcore.AddSync(file) // 输出文件
	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder // 修改时间戳的格式
	encoder := zapcore.NewConsoleEncoder(encoderConfig)
	core := zapcore.NewCore(encoder, writer, zapcore.InfoLevel)
	logT := zap.New(core, zap.AddCaller())
	logger.InfoLogger = logT.Sugar()

	//配置error logger
	file, _ = os.Create("./log/error.txt")
	writer = zapcore.AddSync(file)
	encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	core = zapcore.NewCore(encoder, writer, zapcore.InfoLevel)
	logT = zap.New(core, zap.AddCaller())
	logger.ErrorLogger = logT.Sugar()

	//配置fatal logger
	file, _ = os.Create("./log/fatal.txt")
	writer = zapcore.AddSync(file)
	encoderConfig = zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = zapcore.ISO8601TimeEncoder
	encoder = zapcore.NewConsoleEncoder(encoderConfig)
	core = zapcore.NewCore(encoder, writer, zapcore.InfoLevel)
	logT = zap.New(core, zap.AddCaller())
	logger.FatalLogger = logT.Sugar()
}

func InfoPrintln(str string, userId uint) {
	var printStr string
	if userId != 0 {
		printStr = fmt.Sprintf(str+"\tcause by userId=%d", userId)
	} else {
		printStr = fmt.Sprintf(str + "\tcause by not login user")
	}
	logger.InfoLogger.Infof(printStr)
}

func ErrorPrintln(str string, userId uint, stackByte []byte) {
	printStr := str
	if userId != 0 {
		printStr = fmt.Sprintf(str+"\tcause by userId=%d", userId)
	} else {
		printStr = fmt.Sprintf(str)
	}
	printStr += "\n" + string(stackByte)
	logger.ErrorLogger.Errorf(printStr)
}

func FatalPrintln(str string, userId uint, stackByte []byte) {
	printStr := str
	if userId != 0 {
		printStr = fmt.Sprintf(str+"\tcause by userId=%d", userId)
	} else {
		printStr = fmt.Sprintf(str)
	}
	printStr += "\n" + string(stackByte)
	logger.ErrorLogger.Errorf(printStr)
}

func InitLog() {
	Init()
}
