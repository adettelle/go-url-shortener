package logger

import (
	"log"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func init() {
	var err error
	Logger, err = zap.NewDevelopment()
	if err != nil {
		log.Fatalf("cannot initialize zap: %v", err)
	}
}
