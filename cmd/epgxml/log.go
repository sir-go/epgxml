package main

import (
	"os"

	"github.com/op/go-logging"
)

const logFormat = `%{time:15:04:05.000} %{level:.1s}| %{message}`

var LOG *logging.Logger

func initLogging(appName string) {
	LOG = logging.MustGetLogger(appName)
	formatter := logging.MustStringFormatter(logFormat)
	lb := logging.NewLogBackend(os.Stdout, "", 0)
	lbf := logging.NewBackendFormatter(lb, formatter)
	lbl := logging.AddModuleLevel(lbf)
	logging.SetBackend(lbl)
}
