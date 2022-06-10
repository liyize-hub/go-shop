package utils

import (
	"io"
	"log"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

var Logger *zap.Logger
var SugarLogger *zap.SugaredLogger

//设置infolog, errlog路径 日志级别
func initLog(logPath, errPath string, logLevel zapcore.Level) {
	config := zap.NewProductionEncoderConfig() //编码器(如何写入日志)
	config.EncodeTime = func(t time.Time, pae zapcore.PrimitiveArrayEncoder) {
		pae.AppendString(t.Format("2006-01-02 15:04:05"))
	} // 时间格式
	config.EncodeLevel = zapcore.CapitalLevelEncoder //大写INFO

	// 自定义Info日志级别
	infoLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l < zapcore.WarnLevel && l >= logLevel
	})

	// 自定义Warn日志级别
	warnLevel := zap.LevelEnablerFunc(func(l zapcore.Level) bool {
		return l >= zapcore.WarnLevel && l >= logLevel
	})

	// 获取io.Writer的实现
	infoWriter := getWriter(logPath)
	warnWriter := getWriter(errPath)

	//实现多个输出
	core := zapcore.NewTee(
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(infoWriter), infoLevel),
		zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(warnWriter), warnLevel),
		//zapcore.NewCore(zapcore.NewConsoleEncoder(config), zapcore.AddSync(os.Stdout), logLevel), //打印到命令行界面
	)
	//zap.AddStacktrace(zap.WarnLevel) 错误日志输出调用堆栈
	Logger = zap.New(core, zap.AddCaller()) //加上地址
	SugarLogger = Logger.Sugar()
	Logger.Info("日志启动成功")
}

// 分割日志
func getWriter(filename string) io.Writer {
	// 生产rotatelogs的 logger 实际生成名称为filename.mmddHH
	// filename是指向最新日志的链接
	hook, err := rotatelogs.New(
		filename+".%m%d%H.log",
		rotatelogs.WithLinkName(filename),                          //为最新的日志建立软连接
		rotatelogs.WithRotationTime(time.Duration(60)*time.Second), //切割频率 1h
		rotatelogs.WithMaxAge(time.Duration(1*24)*time.Hour),       //保存时间为 1days
	)
	if err != nil {
		log.Println("日志启动异常")
		log.Fatalln(err)
	}

	return hook
}

func init() {
	initLog("./log/info/info", "./log/err/err", zap.InfoLevel)
}
