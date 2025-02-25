package helper

import (
	"time"
	"os"
	"io"
	"path/filepath"
	"github.com/sirupsen/logrus"
	"github.com/t-tomalak/logrus-easy-formatter"
)

func Logger(logType string, message string) {
	
	time := time.Now().Format("2006-01-02")

	file, _ := os.OpenFile(filepath.Clean("logs/go-"+time+".log"), os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0600)

	logger := &logrus.Logger{
        Out:   os.Stderr,
        Level: logrus.DebugLevel,
        Formatter: &easy.Formatter{
            TimestampFormat: "2006-01-02 15:04:05",
            LogFormat:       "[%lvl%]: %time% - %msg%",
        },
	}

	logger.SetOutput(io.MultiWriter(file, os.Stdout))

	switch logType {
		case "info":
			logger.Info(message + "\n") 
		case "error":
			logger.Error(message + "\n")
	}
}