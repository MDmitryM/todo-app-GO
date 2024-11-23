package main

import (
	"os"

	"github.com/MDmitryM/todo-app-GO"
	"github.com/MDmitryM/todo-app-GO/pkg/handler"
	"github.com/MDmitryM/todo-app-GO/pkg/repository"
	"github.com/MDmitryM/todo-app-GO/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("can't read cfg! %s", err.Error())
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("can't read .env! %s", err.Error())
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("can't open db! %s", err.Error())
	}

	repository := repository.NewRepository(db)
	services := service.NewService(repository)
	handler := handler.NewHandler(services)

	srv := new(todo.Server)

	if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
		logrus.Fatalf("server run error! %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
