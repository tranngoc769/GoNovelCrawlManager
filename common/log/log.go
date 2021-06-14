package log

import (
	"fmt"

	log "github.com/sirupsen/logrus"
)

func Info(origin, function, msg interface{}) {
	log.Info(fmt.Sprintf("%s - %s => msg : %s", origin, function, msg))
}

func Warning(origin, function, msg interface{}) {
	log.Warning(fmt.Sprintf("%s - %s => warn : %s", origin, function, msg))
}

func Error(origin, function string, err interface{}) {
	log.Error(fmt.Sprintf("%s - %s => error : %v", origin, function, err))
}

func Debug(origin, function string, value interface{}) {
	log.Debug(fmt.Sprintf("%s - %s => Debug : %v", origin, function, value))
}

func Fatal(origin, function string, value interface{}) {
	log.Fatal(fmt.Sprintf("%s - %s => Fatal : %v", origin, function, value))
}
