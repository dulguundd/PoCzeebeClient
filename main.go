package main

import (
	"github.com/dulguundd/logError-lib/logger"
	"zeebeClient/application"
)

func main() {
	logger.Info("Starting the application.....")
	application.Start()
	logger.Info("Running")
}
