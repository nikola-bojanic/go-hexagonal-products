package app

import (
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/config"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/database"
	"github.com/mitrovicsoftcoder/go-hexagonal-framework/internal/log"
)

type App struct {
	Config config.Config
	Logger log.Logger

	DB *database.DB
	// AMQP     *amqp.Client
}
