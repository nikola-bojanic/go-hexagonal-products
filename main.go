// Main egw entry point: initializes and starts the server
package main

import (
	"fmt"

	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/app"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/config"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/server"
	"github.com/pkg/errors"
)

func runServer() error {
	app := app.MustInitializeApp()

	cfg := config.ServerConfig{
		Port:   app.Config.Http.Port,
		Logger: app.Logger,
	}

	srv := server.NewServer(cfg, app.DB)

	app.Logger.Infof("Server started at %d", cfg.Port)
	err := srv.ListenAndServe("local", "domain")
	if err != nil {
		return errors.Wrap(err, "listen and serve")
	}
	return nil
}

func main() {
	err := runServer()
	if err != nil {
		fmt.Printf("failed starting app: %s", err)
		panic(err)
	}
}
