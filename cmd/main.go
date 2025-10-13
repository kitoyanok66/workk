package main

import (
	"log"

	"github.com/kitoyanok66/workk/config"
	"github.com/kitoyanok66/workk/internal/db"
	"github.com/kitoyanok66/workk/internal/users"
	"github.com/kitoyanok66/workk/internal/web"
	"github.com/kitoyanok66/workk/internal/web/ousers"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	database, err := db.InitDB(cfg)
	if err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	userRepository := users.NewUserRepository(database)
	userService := users.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)

	strictHandler := ousers.NewStrictHandler(userHandler, nil)

	e := echo.New()
	ousers.RegisterHandlers(e, strictHandler)

}
