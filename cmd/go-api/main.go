package main

import (
	"context"
	"github.com/alpakih/go-api/internal/domain"
	_userHttpDelivery "github.com/alpakih/go-api/internal/users/delivery/http"
	_userRepo "github.com/alpakih/go-api/internal/users/repository/mysql"
	_userService "github.com/alpakih/go-api/internal/users/service"
	"github.com/alpakih/go-api/pkg/database"
	_ "github.com/alpakih/go-api/pkg/database/dialect/postgres"
	"github.com/alpakih/go-api/pkg/env"
	"github.com/alpakih/go-api/pkg/intercept"

	//"github.com/alpakih/go-api/pkg/intercept"

	_ "github.com/alpakih/go-api/pkg/intercept"

	"github.com/alpakih/go-api/pkg/logging"
	"github.com/alpakih/go-api/pkg/validation"
	"github.com/go-oauth2/oauth2/v4"
	"github.com/go-oauth2/oauth2/v4/manage"
	"github.com/go-oauth2/oauth2/v4/models"
	"github.com/go-oauth2/oauth2/v4/server"
	"github.com/go-oauth2/oauth2/v4/store"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"time"
)

var (
	oauthServ *server.Server
	once      sync.Once
)

func main() {

	env.LoadEnvironment()

	db := database.GetConnection()

	manager := manage.NewDefaultManager()
	// token memory store
	manager.MustTokenStorage(store.NewMemoryTokenStore())
	// client memory store
	clientStore := store.NewClientStore()
	clientStore.Set("000000", &models.Client{
		ID:     "000000",
		Secret: "999999",
		Domain: "http://localhost:8080",
	})
	manager.MapClientStorage(clientStore)
	InitServer(manager)
	oauthServ.SetAllowGetAccessRequest(true)
	oauthServ.SetClientInfoHandler(server.ClientFormHandler)

	// Echo instance
	e := echo.New()

	if viper.GetBool("database.autoMigrate") {
		database.RegisterModel(domain.User{})
		database.RegisterModel(domain.OauthClient{})
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

	auth := e.Group("/oauth2")
	{
		auth.POST("/token", HandleTokenRequest)
	}

	// oauth2 password credentials.
	oauthServ.SetPasswordAuthorizationHandler(func(ctx context.Context, username, password string) (userID string, err error) {
		userRepository := _userRepo.NewMysqlUserRepository(db)
		userService := _userService.NewUserService(userRepository)
		result, err := userService.GetByUsername(username)
		if err := bcrypt.CompareHashAndPassword([]byte(result.Password), []byte(password)); err == nil {
			userID = username
		}
		return
	})

	apiGroup := e.Group("/api")
	{
		v1 := apiGroup.Group("/v1")
		{
			userProtect := v1.Group("/users", middleware.JWTWithConfig(intercept.JwtMiddleware().JwtConfig()))
			{
				userRepository := _userRepo.NewMysqlUserRepository(db)
				userService := _userService.NewUserService(userRepository)
				userHandler := _userHttpDelivery.NewUserHandler(userService)
				userProtect.GET("/", userHandler.FetchUsers)
				userProtect.GET("/:id", userHandler.GetUserByID)
				userProtect.PUT("/update", userHandler.UpdateUser)
				userProtect.DELETE("/:id", userHandler.DeleteUser)
			}

			userEndpoint := v1.Group("/users")
			{
				userRepository := _userRepo.NewMysqlUserRepository(db)
				userService := _userService.NewUserService(userRepository)
				userHandler := _userHttpDelivery.NewUserHandler(userService)

				userEndpoint.POST("/token", userHandler.RequestToken)
				userEndpoint.POST("/create", userHandler.StoreUser)
			}

			protectOauth := v1.Group("/bajingan", intercept.OdkOauthCfg(oauthServ, &intercept.DefaultConfig))
			{
				userRepository := _userRepo.NewMysqlUserRepository(db)
				userService := _userService.NewUserService(userRepository)
				userHandler := _userHttpDelivery.NewUserHandler(userService)
				protectOauth.GET("/:id", userHandler.GetUserByID)
			}
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

// InitServer Initialize the service
func InitServer(manager oauth2.Manager) *server.Server {
	once.Do(func() {
		oauthServ = server.NewDefaultServer(manager)
	})
	return oauthServ
}

// HandleTokenRequest token request handling
func HandleTokenRequest(c echo.Context) error {
	return oauthServ.HandleTokenRequest(c.Response().Writer, c.Request())
}

// Sample
func HandleValidate(c echo.Context) error {
	token, err := oauthServ.ValidationBearerToken(c.Request())
	if err != nil {
		return c.JSON(http.StatusBadRequest, echo.Map{"message": http.StatusText(http.StatusBadRequest)})
	}
	return c.JSON(http.StatusOK, echo.Map{"message": token.GetAccess()})
}
