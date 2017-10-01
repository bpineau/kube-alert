package log

import (
	"io"
	"os"

	"log/syslog"

	"github.com/sirupsen/logrus"
	ls "github.com/sirupsen/logrus/hooks/syslog"
)

func New(logLevel string, logServer string, logOutput string) *logrus.Logger {
	var level logrus.Level
	var output io.Writer
	var hook logrus.Hook

	switch logOutput {
	case "stdout":
		output = os.Stdout
	case "stderr":
		output = os.Stderr
	case "syslog":
		output = os.Stderr // does not matter ?
		if logServer == "" {
			panic("syslog output needs a log server (ie. 127.0.0.1:514)")
		}
		hook, _ = ls.NewSyslogHook("udp", logServer, syslog.LOG_INFO, "kube-alert")
	default:
		output = os.Stderr
	}

	switch logLevel {
	case "debug":
		level = logrus.DebugLevel
	case "info":
		level = logrus.InfoLevel
	case "warning":
		level = logrus.WarnLevel
	case "error":
		level = logrus.ErrorLevel
	case "fatal":
		level = logrus.FatalLevel
	case "panic":
		level = logrus.PanicLevel
	default:
		level = logrus.InfoLevel
	}

	log := &logrus.Logger{
		Out:       output,
		Formatter: new(logrus.TextFormatter), //new(logrus.JSONFormatter),
		Hooks:     make(logrus.LevelHooks),
		Level:     level,
	}

	if logOutput == "syslog" {
		log.Hooks.Add(hook)
	}

	return log
}
