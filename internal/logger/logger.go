package logger

import (
	"fmt"
	"go.uber.org/zap"
	"os"
)

var Log *zap.SugaredLogger

func Initialize() error {
	logger, err := zap.NewDevelopment()
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	Log = logger.Sugar()

	return nil
}
