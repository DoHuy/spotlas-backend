package main

import (
	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"spotlas/config"
	"spotlas/http"
	"spotlas/repository"
)

func main() {
	r := gin.Default()
	configuration, err := config.LoadConfigFromEnv()
	if err != nil {
		panic(err)
	}
	db, err := gorm.Open(postgres.Open(configuration.DatabaseURL), &gorm.Config{})
	if err != nil {
		panic(err)
	}

	repo := repository.NewSpotsRepository(db)
	service := http.NewService(repo)
	controller := http.NewController(service)
	http.NewRouter(r, controller).Init()
	if err := r.Run(":3000"); err != nil {
		panic(err)
	}
}
