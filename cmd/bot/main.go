package main

import (
	"log"
	"log/slog"
	"os"

	"github.com/Jacobbrewer1/wolf/pkg/logging"
)

func main() {
	a, err := InitializeApp()
	if err != nil {
		log.Fatalln(err)
	}
	parseConfig()
	a.Info("Starting application")
	if err := a.Run(); err != nil {
		slog.Error("Error running application", slog.String(logging.KeyError, err.Error()))
		os.Exit(1)
	}
}
