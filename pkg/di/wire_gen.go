// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package di

import (
	"video/pkg/api"
	"video/pkg/api/service"
	"video/pkg/config"
	"video/pkg/db"
	"video/pkg/repo"
)

// Injectors from wire.go:

func InitializeServer(c *config.Config) (*api.Server, error) {
	gormDB, err := db.InitDB(c)
	if err != nil {
		return nil, err
	}
	videoRepo := repo.NewVideoRepo(gormDB)
	videoServiceServer := service.NewVideoServer(videoRepo)
	server, err := api.NewGrpcServer(c, videoServiceServer)
	if err != nil {
		return nil, err
	}
	return server, nil
}
