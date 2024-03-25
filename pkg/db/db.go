package db

import (
	"fmt"
	"log"
	"video/pkg/config"
	"video/pkg/models"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func InitDB(cfg *config.Config) (*gorm.DB,error) {
	psql:=fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s", cfg.Db_host, cfg.Db_username, cfg.Db_password, cfg.Db_name, cfg.Db_port)
	db,err:=gorm.Open(postgres.Open(psql),&gorm.Config{})
	if err != nil{
		log.Fatalln(err)
	}
	db.AutoMigrate(
		&models.Video{},
		&models.Viewer{},
	)
	return db,err
	}	