//go:build wireinject
// +build wireinject

package app

import (
	"github.com/google/wire"
	"github.com/kitoyanok66/workk/config"
	"github.com/kitoyanok66/workk/internal/db"
	"github.com/kitoyanok66/workk/internal/freelancers"
	"github.com/kitoyanok66/workk/internal/likes"
	"github.com/kitoyanok66/workk/internal/matches"
	"github.com/kitoyanok66/workk/internal/projects"
	"github.com/kitoyanok66/workk/internal/skills"
	"github.com/kitoyanok66/workk/internal/users"
	"gorm.io/gorm"
)

type App struct {
	DB                *gorm.DB
	UserService       users.UserService
	SkillService      skills.SkillService
	ProjectService    projects.ProjectService
	FreelancerService freelancers.FreelancerService
	LikeService       likes.LikeService
	MatchService      matches.MatchService
}

func InitApp(cfg *config.Config) (*App, error) {
	wire.Build(
		db.InitDB,

		users.NewUserRepository,
		skills.NewSkillRepository,
		likes.NewLikeRepository,
		matches.NewMatchRepository,
		freelancers.NewFreelancerRepository,
		projects.NewProjectRepository,

		users.NewUserService,
		skills.NewSkillService,
		likes.NewLikeService,
		matches.NewMatchService,
		freelancers.NewFreelancerService,
		projects.NewProjectService,

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
	}
}
