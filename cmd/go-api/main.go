package main

import (
	"context"
	"github.com/alpakih/go-api/internal/domain"
	_userHttpDelivery "github.com/alpakih/go-api/internal/users/delivery/http"
	_userRepo "github.com/alpakih/go-api/internal/users/repository/mysql"
	_userService "github.com/alpakih/go-api/internal/users/service"
	"github.com/alpakih/go-api/pkg/database"
	_ "github.com/alpakih/go-api/pkg/database/dialect/mysql"
	"github.com/alpakih/go-api/pkg/env"
	"github.com/alpakih/go-api/pkg/intercept"
	"github.com/alpakih/go-api/pkg/logging"
	"github.com/alpakih/go-api/pkg/validation"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"
)

func main() {

	env.LoadEnvironment()

	// Echo instance
	e := echo.New()

	db := database.GetConnection()

	if viper.GetBool("database.autoMigrate") {
		database.RegisterModel(domain.User{})
		database.Migrate()
	}

	// Set Middleware
	e.Use(middleware.Recover())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAccept, echo.HeaderAuthorization},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))

	e.Validator = validation.NewValidator()

	// setup log folder,file and log global
	logFile := logging.SetupLogFileAndFolder("api")

	// set logger middleware
	e.Use(middleware.LoggerWithConfig(middleware.LoggerConfig{
		Output: io.MultiWriter(logFile, os.Stdout),
	}))

	// default handler
	e.GET("/", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "GO API")
	})

	apiGroup := e.Group("/api")
	{
		v1 := apiGroup.Group("/v1", middleware.JWTWithConfig(intercept.JwtMiddleware().JwtConfig()))
		{
			userRepository := _userRepo.NewMysqlUserRepository(db)
			userService := _userService.NewUserService(userRepository)
			userHandler := _userHttpDelivery.NewUserHandler(userService)

			v1.POST("/users/token", userHandler.RequestToken)
			v1.GET("/users", userHandler.FetchUsers)
			v1.GET("/users", userHandler.FetchUsers)
			v1.GET("/users/:id", userHandler.GetUserByID)
			v1.POST("/users", userHandler.StoreUser)
			v1.PUT("/users", userHandler.UpdateUser)
			v1.DELETE("/users/:id", userHandler.DeleteUser)
		}
	}

	// Start server
	go func() {
		if err := e.Start(viper.GetString("server.host") + ":" + viper.GetString("server.port")); err != nil {
			e.Logger.Info("shutting down the server")
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with a timeout of 10 seconds.
	// Use a buffered channel to avoid missing signals as recommended for signal.Notify
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err := e.Shutdown(ctx); err != nil {
		e.Logger.Fatal(err)
	}
}
