package http

import (
	"errors"
	"github.com/alpakih/go-api/internal/domain"
	"github.com/alpakih/go-api/pkg/validation"
	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/gommon/log"
	"github.com/spf13/viper"
	"golang.org/x/crypto/bcrypt"
	"gopkg.in/go-playground/validator.v9"
	"gorm.io/gorm"
	"net/http"
	"strconv"
	"time"
)

type UserHandler struct {
	UserService domain.UserService
}

func NewUserHandler(us domain.UserService) UserHandler {
	return UserHandler{UserService: us}
}

func (r *UserHandler) RequestToken(ctx echo.Context) error {
	var request domain.TokenRequest
	var accessToken string
	if err := ctx.Bind(&request); err != nil {
		log.Error(err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": http.StatusText(http.StatusBadRequest)})
	}
	if err := ctx.Validate(&request); err != nil {
		log.Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity,
			echo.Map{"message": http.StatusText(http.StatusUnprocessableEntity),
				"errors": validation.WrapValidationErrors(err.(validator.ValidationErrors))})
	}

	result, err := r.UserService.GetByUsername(request.Username)
	if err != nil {
		log.Error(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, echo.Map{"message": http.StatusText(http.StatusNotFound)})
		}
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": http.StatusText(http.StatusInternalServerError)})
	}

	if err := bcrypt.CompareHashAndPassword( []byte(result.Password),[]byte(request.Password)); err != nil {
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": "incorrect username or password"})
	}

	// Create token with claims
	token := jwt.New(jwt.SigningMethodHS256)
	tokenClaims := token.Claims.(jwt.MapClaims)

	tokenClaims["id"] = result.ID
	tokenClaims["username"] = result.UserName
	tokenClaims["exp"] = time.Now().Add(time.Duration(viper.GetInt("auth.jwt.validity")) * time.Second).Unix()

	//Encode Token
	if encodeToken, err := token.SignedString([]byte(viper.GetString("auth.jwt.secret"))); err != nil {
		log.Error(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": http.StatusText(http.StatusInternalServerError)})
	} else {
		accessToken = encodeToken
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": http.StatusText(http.StatusOK),
		"data": echo.Map{"access_token": accessToken, "exp": tokenClaims["exp"]}})
}

func (r *UserHandler) FetchUsers(ctx echo.Context) error {
	params := ctx.QueryParams()

	// default vars
	limit := 10
	offset := 0

	if offsetParse, err := strconv.Atoi(params.Get("offset")); err == nil {
		offset = offsetParse
	}
	if limitParse, err := strconv.Atoi(params.Get("limit")); err == nil {
		limit = limitParse
	}

	result, err := r.UserService.Fetch(limit, offset)
	if err != nil {
		log.Error(err)
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": http.StatusText(http.StatusInternalServerError)})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"message": http.StatusText(http.StatusOK), "data": result})
}

func (r *UserHandler) GetUserByID(ctx echo.Context) error {
	param := ctx.Param("id")

	result, err := r.UserService.GetByID(param)
	if err != nil {
		log.Error(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, echo.Map{"message": http.StatusText(http.StatusNotFound)})
		}
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": http.StatusText(http.StatusInternalServerError)})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"message": http.StatusText(http.StatusOK), "data": result})
}

func (r *UserHandler) StoreUser(ctx echo.Context) error {

	var request domain.StoreRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": http.StatusText(http.StatusBadRequest)})
	}
	if err := ctx.Validate(&request); err != nil {
		log.Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity,
			echo.Map{"message": http.StatusText(http.StatusUnprocessableEntity),
				"errors": validation.WrapValidationErrors(err.(validator.ValidationErrors))})
	}

	if err := r.UserService.Store(request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": http.StatusText(http.StatusInternalServerError)})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "save data success"})
}

func (r *UserHandler) UpdateUser(ctx echo.Context) error {

	var request domain.UpdateRequest
	if err := ctx.Bind(&request); err != nil {
		log.Error(err)
		return ctx.JSON(http.StatusBadRequest, echo.Map{"message": http.StatusText(http.StatusBadRequest)})
	}
	if err := ctx.Validate(&request); err != nil {
		log.Error(err)
		return ctx.JSON(http.StatusUnprocessableEntity,
			echo.Map{"message": http.StatusText(http.StatusUnprocessableEntity),
				"errors": validation.WrapValidationErrors(err.(validator.ValidationErrors))})
	}

	if _, err := r.UserService.GetByID(request.ID); err != nil {
		log.Error(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, echo.Map{"message": http.StatusText(http.StatusNotFound)})
		}
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": http.StatusText(http.StatusInternalServerError)})
	}

	if err := r.UserService.Update(request); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": http.StatusText(http.StatusInternalServerError)})
	}

	return ctx.JSON(http.StatusOK, echo.Map{"message": "update data success"})
}

func (r *UserHandler) DeleteUser(ctx echo.Context) error {
	param := ctx.Param("id")

	if _, err := r.UserService.GetByID(param); err != nil {
		log.Error(err)
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ctx.JSON(http.StatusNotFound, echo.Map{"message": http.StatusText(http.StatusNotFound)})
		}
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": http.StatusText(http.StatusInternalServerError)})
	}
	if err := r.UserService.Delete(param); err != nil {
		return ctx.JSON(http.StatusInternalServerError, echo.Map{"message": http.StatusText(http.StatusInternalServerError)})
	}
	return ctx.JSON(http.StatusOK, echo.Map{"message": "delete data success"})
}
