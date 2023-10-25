package logger

import (
	"os"

	"github.com/sirupsen/logrus"
)

var Log = logrus.New()

func Init() {
	Log.Out = os.Stdout
	Log.Level = logrus.InfoLevel
	Log.Formatter = &logrus.TextFormatter{
		DisableTimestamp: true,
	}
}
