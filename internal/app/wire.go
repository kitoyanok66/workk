//go:build wireinject
// +build wireinject

package app

import (
	"time"

	"github.com/google/wire"
	"github.com/kitoyanok66/workk/config"
	"github.com/kitoyanok66/workk/internal/auth"
	"github.com/kitoyanok66/workk/internal/db"
	"github.com/kitoyanok66/workk/internal/freelancers"
	"github.com/kitoyanok66/workk/internal/likes"
	"github.com/kitoyanok66/workk/internal/matches"
	"github.com/kitoyanok66/workk/internal/middleware"
	"github.com/kitoyanok66/workk/internal/projects"
	"github.com/kitoyanok66/workk/internal/skills"
	"github.com/kitoyanok66/workk/internal/users"
	"gorm.io/gorm"
)

type App struct {
	DB                *gorm.DB
	JWTManager        *middleware.JWTManager
	UserService       users.UserService
	SkillService      skills.SkillService
	ProjectService    projects.ProjectService
	FreelancerService freelancers.FreelancerService
	LikeService       likes.LikeService
	MatchService      matches.MatchService
	AuthService       auth.AuthService
}

func provideJWTSecret(cfg *config.Config) string {
	return cfg.JWTSecret
}

func provideJWTTTL(cfg *config.Config) time.Duration {
	return cfg.JWTTTLHours
}

func InitApp(cfg *config.Config) (*App, error) {
	wire.Build(
		db.InitDB,

		provideJWTSecret,
		provideJWTTTL,
		middleware.NewJWTManager,

		users.NewUserRepository,
		skills.NewSkillRepository,
		likes.NewLikeRepository,
		matches.NewMatchRepository,
		freelancers.NewFreelancerRepository,
		projects.NewProjectRepository,
		auth.NewAuthRepository,

		users.NewUserService,
		skills.NewSkillService,
		likes.NewLikeService,
		matches.NewMatchService,
		freelancers.NewFreelancerService,
		projects.NewProjectService,
		auth.NewAuthService,

		NewApp,
	)
	return nil, nil
}

func NewApp(
	db *gorm.DB,
	userSvc users.UserService,
	skillSvc skills.SkillService,
	projectSvc projects.ProjectService,
	freelancerSvc freelancers.FreelancerService,
	likeSvc likes.LikeService,
	matchSvc matches.MatchService,
	authSvc auth.AuthService,
	jwtManager *middleware.JWTManager,
) *App {
	fAdapter := freelancers.NewFreelancerFetcherAdapter(freelancerSvc)
	pAdapter := projects.NewProjectFetcherAdapter(projectSvc)

	freelancerSvc.SetProjectDep(pAdapter)
	projectSvc.SetFreelancerDep(fAdapter)

	return &App{
		DB:                db,
		UserService:       userSvc,
		SkillService:      skillSvc,
		ProjectService:    projectSvc,
		FreelancerService: freelancerSvc,
		LikeService:       likeSvc,
		MatchService:      matchSvc,
		AuthService:       authSvc,
		JWTManager:        jwtManager,
	}
}
