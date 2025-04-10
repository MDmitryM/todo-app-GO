package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/MDmitryM/todo-app-GO"
	"github.com/MDmitryM/todo-app-GO/pkg/handler"
	"github.com/MDmitryM/todo-app-GO/pkg/repository"
	"github.com/MDmitryM/todo-app-GO/pkg/service"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter))
	if err := initConfig(); err != nil {
		logrus.Fatalf("can't read cfg! %s", err.Error())
	}

	var conf repository.Config
	env := os.Getenv("ENV")
	if env != "production" {
		// Загружаем .env только в dev-окружении
		if err := godotenv.Load(); err != nil {
			logrus.Fatalf("No .env file found, relying on environment variables")
		}
		conf = repository.Config{
			Host:     viper.GetString("dev_db.host"),
			Port:     viper.GetString("dev_db.port"),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
			SSLMode:  "disable",
		}
		logrus.Println("Using localhost for DB connection")
	} else {
		conf = repository.Config{
			Host:     viper.GetString("db.host"),
			Port:     viper.GetString("db.port"),
			Username: os.Getenv("POSTGRES_USER"),
			Password: os.Getenv("DB_PASSWORD"),
			DBName:   os.Getenv("POSTGRES_DB"),
			SSLMode:  "disable",
		}
	}
	fmt.Printf("k=%s s=%s", os.Getenv("SIGNING_KEY"), os.Getenv("SALT"))
	db, err := repository.NewPostgresDB(conf)
	if err != nil {
		logrus.Fatalf("can't open db! %s", err.Error())
	}

	repository := repository.NewRepository(db)
	services, err := service.NewService(repository)
	if err != nil {
		logrus.Fatalf("create service error! %s", err.Error())
	}
	handler := handler.NewHandler(services)

	srv := new(todo.Server)

	go func() {
		if err := srv.Run(viper.GetString("port"), handler.InitRoutes()); err != nil {
			if err != http.ErrServerClosed {
				logrus.Fatalf("server run error! %s", err.Error())
			}
		}
	}()

	logrus.Print("Server started, todo app is running.")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("shutting down")

	if err := srv.Shutdown(context.Background()); err != nil {
		logrus.Errorf("error occured on server shutting down: %s", err.Error())
	}

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on database close: %s", err.Error())
	}
}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
