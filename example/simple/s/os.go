package s

import (
	"os"
	"os/signal"
	"syscall"

	"github.com/vault-thirteen/auxie/cio"
)

func WaitForQuitSignalFromOS() {
	osSignals := make(chan os.Signal, 1)
	signal.Notify(osSignals, syscall.SIGINT, syscall.SIGTERM)

	for sig := range osSignals {
		switch sig {
		case syscall.SIGINT,
			syscall.SIGTERM:
			cio.LogMessage(sig)
			return
		}
	}
}
