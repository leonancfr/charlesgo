package gablogger

import (
	"common"
	"fmt"
	"path"
	"runtime"
	"time"

	datadog "github.com/bin3377/logrus-datadog-hook"
	log "github.com/sirupsen/logrus"
)

var myLogger = log.New()

func ConfigureDatadog(hostName string) error {
	if common.DATADOG_HOST == "" || common.DATADOG_API_KEY == "" {
		return fmt.Errorf("datadog credentials not defined")
	}

	hook := datadog.NewHook(
		common.DATADOG_HOST,
		common.DATADOG_API_KEY,
		time.Minute, // Batch timeout
		3,           // MaxRetry
		log.InfoLevel,
		&log.JSONFormatter{},
		datadog.Options{
			Hostname: hostName,
			Service:  "CharlesGo",
			Source:   common.ENVIRONMENT,
			Tags:     []string{},
		},
	)
	myLogger.Hooks.Add(hook)
	return nil
}

func init() {
	logLevel := log.DebugLevel
	if common.ENVIRONMENT == "production" {
		logLevel = log.InfoLevel
	}

	myLogger.SetLevel(logLevel)
	myLogger.SetReportCaller(true)
	myLogger.SetFormatter(&log.TextFormatter{
		FullTimestamp:   true,
		TimestampFormat: "2006-01-02 15:04:05",
		CallerPrettyfier: func(f *runtime.Frame) (string, string) {
			fileName := path.Base(f.File)
			return "", fmt.Sprintf("[%s:%d]", fileName, f.Line)
		},
	})

}

func Logger() *log.Logger {
	return myLogger
}
