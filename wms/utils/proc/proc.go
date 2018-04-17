package proc

import (
	"fmt"
	"message_server/utils/file"
	"message_server/utils/logs"
	"os"
	"os/signal"
	"syscall"
)

var (
	FileName = "run.pid"
)

func init() {
	start()
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-signalChan
		stop()
		os.Exit(0)
	}()
}

func start() {
	err := file.Writer("./"+FileName, fmt.Sprintf("%d", os.Getpid()))
	if err != nil {
		logs.Error(err.Error())
	}
}

func stop() {
	err := file.Remove("./" + FileName)
	if err != nil {
		logs.Error(err.Error())
	}
}
