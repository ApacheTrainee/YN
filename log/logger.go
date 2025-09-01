package log

import (
	"YN/config"
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"github.com/sirupsen/logrus"
)

var Logger *logrus.Logger
var WebLogger *logrus.Logger

// 自定义logrus输出格式
type myFormatter struct{}

func (m *myFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	var b *bytes.Buffer
	if entry.Buffer != nil {
		b = entry.Buffer
	} else {
		b = &bytes.Buffer{}
	}

	var newLog string
	timestamp := entry.Time.Format("2006/01/02 - 15:04:05")

	// SetReportCaller(true)即entry.HasCaller()为true
	if entry.HasCaller() && entry.Level >= logrus.InfoLevel {
		newLog = fmt.Sprintf("[%s] %s: %s \n", entry.Level, timestamp, entry.Message)
	} else {
		newLog = fmt.Sprintf("[%s] %s %s:%d: %s \n", entry.Level, timestamp, entry.Caller.File, entry.Caller.Line, entry.Message)
	}

	b.WriteString(newLog)
	return b.Bytes(), nil
}

// 自定义logrus错误输出Hook
type errRotateFileHook struct {
	writer    io.Writer
	formatter logrus.Formatter
}

func (hook *errRotateFileHook) Levels() []logrus.Level { // 捕抓到以下日志等级时，触发Hook
	return []logrus.Level{
		logrus.ErrorLevel,
		logrus.PanicLevel,
		logrus.FatalLevel,
	}
}

func (hook *errRotateFileHook) Fire(entry *logrus.Entry) error {
	serialized, err := hook.formatter.Format(entry)
	if err != nil {
		return fmt.Errorf("failed to format entry: %w", err)
	}

	if _, err = hook.writer.Write(serialized); err != nil {
		return fmt.Errorf("failed to write to error log file: %w", err)
	}
	return nil
}

func logSet(finalDirectory string) *logrus.Logger {
	// () 创建logrus对象，设置配置
	logrusLog := logrus.New()
	logrusLog.SetReportCaller(true)        // 启用Caller(启用文件名、行号记录功能)
	logrusLog.SetFormatter(&myFormatter{}) // 自定义输出格式

	logrusLog.SetLevel(logrus.DebugLevel) // 达到什么日志等级输出
	if config.Config.RunMode == "pro" {
		logrusLog.SetLevel(logrus.InfoLevel)
	}

	// () 定义日志输出目录
	currentDir, errPath := os.Getwd() // 获取当前目录
	if errPath != nil {
		panic("Failed to obtain the absolute path of the current directory" + errPath.Error())
	}

	logOutputPath := filepath.Join(currentDir, "log", finalDirectory) // 定义日志输入路径
	if err := os.MkdirAll(logOutputPath, 0755); err != nil {
		panic("Failed to create log directory" + logOutputPath + err.Error())
	}

	// () 配置Hook：额外输出错误日志
	writerErr, _ := rotatelogs.New(
		filepath.Join(logOutputPath, "%Y-%m-%d", "YN_log_err.%Y-%m-%d_%H.log"), // %Y%m%d%H%M%S: 年月日 时分秒
		rotatelogs.WithMaxAge(30*24*time.Hour),
		//rotatelogs.WithRotationTime(1*time.Hour),
	)
	errFileHook := &errRotateFileHook{writer: writerErr, formatter: &myFormatter{}}
	logrusLog.AddHook(errFileHook)

	// () 配置输出日志
	writer, _ := rotatelogs.New(
		filepath.Join(logOutputPath, "%Y-%m-%d", "YN_log.%Y-%m-%d_%H.log"), // %Y%m%d %H%M%S: 年月日 时分秒
		rotatelogs.WithMaxAge(30*24*time.Hour),
		//rotatelogs.WithRotationTime(1*time.Hour),
	)
	logrusLog.SetOutput(io.MultiWriter(os.Stdout, writer))

	return logrusLog
}

func InitLog() {
	Logger = logSet("logfile")
	WebLogger = logSet("logfile_web")
}
