package main

import (
	"log"
	"os"

	"github.com/Jacobbrewer1/wolf/cmd/bot/config"
	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"golang.org/x/exp/slog"
)

func main() {
	a, err := InitializeApp()
	if err != nil {
		log.Fatalln(err)
	}
	config.Parse(a.Log())
	a.Info("Starting application")
	if err := a.Run(); err != nil {
		a.Error("Error running application", slog.String(logging.KeyError, err.Error()))
		os.Exit(1)
	}
}
