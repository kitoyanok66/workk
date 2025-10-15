package main

import (
	"log"

	"github.com/kitoyanok66/workk/config"
	"github.com/kitoyanok66/workk/internal/db"
	"github.com/kitoyanok66/workk/internal/freelancers"
	"github.com/kitoyanok66/workk/internal/skills"
	"github.com/kitoyanok66/workk/internal/users"
	"github.com/kitoyanok66/workk/internal/web"
	"github.com/kitoyanok66/workk/internal/web/ofreelancers"
	"github.com/kitoyanok66/workk/internal/web/oskills"
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

	e := echo.New()

	userRepository := users.NewUserRepository(database)
	userService := users.NewUserService(userRepository)
	userHandler := web.NewUserHandler(userService)
	strictUserHandler := ousers.NewStrictHandler(userHandler, nil)
	ousers.RegisterHandlers(e, strictUserHandler)

	skillRepository := skills.NewSkillRepository(database)
	skillService := skills.NewSkillService(skillRepository)
	skillHandler := web.NewSkillHandler(skillService)
	strictSkillHandler := oskills.NewStrictHandler(skillHandler, nil)
	oskills.RegisterHandlers(e, strictSkillHandler)

	freelancerRepository := freelancers.NewSkillRepository(database)
	freelancerService := freelancers.NewFreelancerService(freelancerRepository, skillRepository)
	freelancerHandler := web.NewFreelancerHandler(freelancerService)
	strictFreelancerHandler := ofreelancers.NewStrictHandler(freelancerHandler, nil)
	ofreelancers.RegisterHandlers(e, strictFreelancerHandler)
}
