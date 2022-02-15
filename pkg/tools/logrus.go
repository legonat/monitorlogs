package tools

import (
	"fmt"
	logrus "github.com/sirupsen/logrus"
	"io"
	"os"
)

type PlainFormatter struct {
	TimestampFormat string
	LevelDesc       []string
}

func (f *PlainFormatter) Format(entry *logrus.Entry) ([]byte, error) {
	timestamp := fmt.Sprintf(entry.Time.Format(f.TimestampFormat))
	return []byte(fmt.Sprintf("[%s] [%s] %s\n", timestamp, f.LevelDesc[entry.Level], entry.Message)), nil
}

func Logrus(file *os.File) (*logrus.Logger, error) {
	log := logrus.New()
	mw := io.MultiWriter(os.Stdout, file)
	log.SetOutput(mw)
	plainFormatter := new(PlainFormatter)
	plainFormatter.TimestampFormat = "2006-01-02 15:04:05 MST"
	plainFormatter.LevelDesc = []string{"PANIC", "FATAL", "ERROR", "WARNING", "INFO", "DEBUG"}
	log.SetFormatter(plainFormatter)
	//logger := &log. {
	//	Out: mw,
	//	Formatter: &prefixed.TextFormatter{
	//		DisableColors: true,
	//		TimestampFormat : "2006-01-02 15:04:05 MST",
	//		FullTimestamp:true,
	//		ForceFormatting: true,
	//	},
	//}

	//log.SetFormatter(&prefixed.TextFormatter{TimestampFormat: "2006-01-02 15:04:05 MST", FullTimestamp: true, ForceFormatting: true})

	return log, nil
}

func LogInfo(msg string) {
	path := os.Getenv("LOGGER")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		return
	}
	defer file.Close()
	logger, err := Logrus(file)
	if err != nil {
		return
	}
	logger.Info(msg)
}

func LogWarn(msg string) {
	path := os.Getenv("LOGGER")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	defer file.Close()
	if err != nil {
		return
	}
	logger, err := Logrus(file)
	if err != nil {
		return
	}
	logger.Warn(msg)
}

func LogErr(exErr error) {
	path := os.Getenv("LOGGER")
	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	defer file.Close()
	if err != nil {
		return
	}
	logger, err := Logrus(file)
	if err != nil {
		return
	}
	logger.Error(exErr)
}
