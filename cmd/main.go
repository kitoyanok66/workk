package main

import (
	"log"

	"github.com/kitoyanok66/workk/config"
	"github.com/kitoyanok66/workk/internal/app"
	"github.com/kitoyanok66/workk/internal/web"
	"github.com/kitoyanok66/workk/internal/web/ofreelancers"
	"github.com/kitoyanok66/workk/internal/web/oprojects"
	"github.com/kitoyanok66/workk/internal/web/oskills"
	"github.com/kitoyanok66/workk/internal/web/ousers"
	"github.com/labstack/echo/v4"
)

func main() {
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	a, err := app.InitApp(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	e := echo.New()

	userHandler := web.NewUserHandler(a.UserService)
	strictUserHandler := ousers.NewStrictHandler(userHandler, nil)
	ousers.RegisterHandlers(e, strictUserHandler)

	skillHandler := web.NewSkillHandler(a.SkillService)
	strictSkillHandler := oskills.NewStrictHandler(skillHandler, nil)
	oskills.RegisterHandlers(e, strictSkillHandler)

	freelancerHandler := web.NewFreelancerHandler(a.FreelancerService)
	strictFreelancerHandler := ofreelancers.NewStrictHandler(freelancerHandler, nil)
	ofreelancers.RegisterHandlers(e, strictFreelancerHandler)

	projectHandler := web.NewProjectHandler(a.ProjectService)
	strictProjectHandler := oprojects.NewStrictHandler(projectHandler, nil)
	oprojects.RegisterHandlers(e, strictProjectHandler)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with error: %v", err)
	}
}
