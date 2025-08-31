package logging

import (
	"log"

	"go.uber.org/zap"
)

var Logger *zap.Logger

func InitLogger() {
	cfg := zap.NewProductionConfig()
	cfg.OutputPaths = []string{"stdout"}

	var err error
	Logger, err = cfg.Build()
	if err != nil {
		log.Fatalf("failed to initialize zap logger: %v", err)
	}

	// if err != nil {
	// 	log.Fatalf("failed to initialize zap logger: %v", err)
	// }
	// a.logger = logr
}
