//go:build wireinject
// +build wireinject

package main

import (
	"github.com/Jacobbrewer1/wolf/pkg/logging"
	"github.com/google/wire"
	"github.com/gorilla/mux"
)

func InitializeApp() (*App, error) {
	wire.Build(
		wire.Value(logging.Name(AppName)),
		logging.NewConfig,
		logging.CommonLogger,
		mux.NewRouter,
		NewApp,
	)
	return new(App), nil
}
