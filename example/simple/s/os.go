package s

import (
	"log"
	"os"
	"os/signal"
	"syscall"
)

func WaitForQuitSignalFromOS() {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	for sig := range osSignals {
		switch sig {
		case syscall.SIGINT,
			syscall.SIGTERM:
			log.Println(sig)
			return
		}
	}
}
