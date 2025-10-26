package main

import (
	"log"

	"github.com/kitoyanok66/workk/config"
	"github.com/kitoyanok66/workk/internal/app"
	"github.com/kitoyanok66/workk/internal/web"
	"github.com/kitoyanok66/workk/internal/web/oauth"
	"github.com/kitoyanok66/workk/internal/web/ofreelancers"
	"github.com/kitoyanok66/workk/internal/web/olikes"
	"github.com/kitoyanok66/workk/internal/web/omatches"
	"github.com/kitoyanok66/workk/internal/web/oprojects"
	"github.com/kitoyanok66/workk/internal/web/oskills"
	"github.com/kitoyanok66/workk/internal/web/ousers"
	"github.com/labstack/echo/v4"
)

func main() {
	// загружаем конфиг
	cfg, err := config.Load()
	if err != nil {
		log.Fatalf("failed to load config: %v", err)
	}

	// инициализируем приложение через wire
	a, err := app.InitApp(cfg)
	if err != nil {
		log.Fatalf("failed to initialize app: %v", err)
	}

	e := echo.New()

	// публичные маршруты
	authHandler := web.NewAuthHandler(a.AuthService, a.UserService, a.JWTManager)
	strictAuthHandler := oauth.NewStrictHandler(authHandler, nil)
	oauth.RegisterHandlers(e, strictAuthHandler)

	// защищенные маршруты

	// создаем группу для JWT-проверки
	protected := e.Group("")
	protected.Use(a.JWTManager.JWTStrictMiddleware())

	// users
	userHandler := web.NewUserHandler(a.UserService)
	strictUserHandler := ousers.NewStrictHandler(userHandler, nil)
	ousers.RegisterHandlers(protected, strictUserHandler)

	// skills
	skillHandler := web.NewSkillHandler(a.SkillService)
	strictSkillHandler := oskills.NewStrictHandler(skillHandler, nil)
	oskills.RegisterHandlers(protected, strictSkillHandler)

	// freelancers
	freelancerHandler := web.NewFreelancerHandler(a.FreelancerService)
	strictFreelancerHandler := ofreelancers.NewStrictHandler(freelancerHandler, nil)
	ofreelancers.RegisterHandlers(protected, strictFreelancerHandler)

	// projects
	projectHandler := web.NewProjectHandler(a.ProjectService)
	strictProjectHandler := oprojects.NewStrictHandler(projectHandler, nil)
	oprojects.RegisterHandlers(protected, strictProjectHandler)

	// likes
	likeHandler := web.NewLikeHandler(a.LikeService, a.MatchService, a.FreelancerService, a.ProjectService, a.UserService)
	strictLikeHandler := olikes.NewStrictHandler(likeHandler, nil)
	olikes.RegisterHandlers(protected, strictLikeHandler)

	// matches
	matchHandler := web.NewMatchHandler(a.MatchService, a.FreelancerService, a.ProjectService, a.UserService)
	strictMatchHandler := omatches.NewStrictHandler(matchHandler, nil)
	omatches.RegisterHandlers(protected, strictMatchHandler)

	if err := e.Start(":8080"); err != nil {
		log.Fatalf("failed to start with error: %v", err)
	}
}
