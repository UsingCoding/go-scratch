package main

import (
	"time"

	"github.com/sirupsen/logrus"
)

func initLogger() *logrus.Logger {
	impl := logrus.New()
	impl.SetFormatter(&logrus.TextFormatter{
		TimestampFormat:  time.RFC3339Nano,
		DisableTimestamp: true,
	})

	return impl
}
