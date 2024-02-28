package main

import (
	"Proxy/internal/handler"
	"Proxy/internal/proxy"
	"Proxy/internal/repo"
	"Proxy/internal/usecase"
	"github.com/spf13/viper"
	"log"
	"path/filepath"
	"strings"
)

func main() {
	vp := viper.New()
	if err := initConfig(vp, "/configs/config.yml"); err != nil {
		log.Printf("error initializing configs: %s\n", err.Error())
	}

	db, err := repo.NewPostgresDB(vp.GetString("db.connection_string"))
	if err != nil {
		log.Fatal("error during connecting to postgres ", err)
	}

	repos := repo.NewRepo(db)
	services := usecase.NewUsecase(repos)
	handlers := handler.NewHandler(services)

	go func() {
		proxy.Run(services)
	}()

	router := handlers.SetupRoutes()
	router.Run("0.0.0.0:8000")
}

func initConfig(vp *viper.Viper, configPath string) error {
	path := filepath.Dir(configPath)
	vp.AddConfigPath(path)
	vp.SetConfigName(strings.Split(filepath.Base(configPath), ".")[0])
	vp.SetConfigType(filepath.Ext(configPath)[1:])
	return vp.ReadInConfig()
}
