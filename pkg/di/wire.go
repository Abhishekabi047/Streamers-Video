package di

import (
	"video/pkg/api"
	"video/pkg/api/service"
	"video/pkg/config"
	"video/pkg/db"
	"video/pkg/repo"

	"github.com/google/wire"
)

func InitializeServe(c *config.Config) (*api.Server,error) {
	wire.Build(db.InitDB,repo.NewVideoRepo,service.NewVideoServer,api.NewGrpcServer)
	return &api.Server{},nil
}