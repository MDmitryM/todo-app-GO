package main

import (
	"log"

	"github.com/MDmitryM/todo-app-GO"
	"github.com/MDmitryM/todo-app-GO/pkg/handler"
	"github.com/MDmitryM/todo-app-GO/pkg/repository"
	"github.com/MDmitryM/todo-app-GO/pkg/service"
	"github.com/spf13/viper"
)

func main() {
	if err := initConfig(); err != nil {
		log.Fatalf("can't read cfg! %s", err.Error())
	}

	repository := repository.NewRepository()
	services := service.NewService(repository)
	handler := handler.NewHandler(services)

	srv := new(todo.Server)

	if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		log.Fatalf("server run error! %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
