package handler

import (
	"database/sql"
	"encoding/base64"
	"fmt"
	"github.com/SawitProRecruitment/UserService/generated"
	"github.com/SawitProRecruitment/UserService/repository"
	"github.com/labstack/echo/v4"
	"golang.org/x/crypto/bcrypt"
	"net/http"
	"reflect"
	"strings"
	"unicode"
)

func (s *Server) Login(ctx echo.Context) error {
	var request generated.LoginRequest
	if err := ctx.Bind(&request); nil != err {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "invalid parameter request"})
	}

	if err := validatePhone(request.Phone); nil != err {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	users, err := s.Repository.FindByPhone(ctx.Request().Context(), repository.FindByPhoneInput{Phone: request.Phone})
	if nil != err {
		if err == sql.ErrNoRows {
			return ctx.JSON(http.StatusNotFound, generated.ErrorResponse{Message: "user not found"})
		}

		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: http.StatusText(http.StatusInternalServerError)})
	}

	if err = bcrypt.CompareHashAndPassword([]byte(users.Password), []byte(request.Password)); nil != err {
		return ctx.JSON(http.StatusForbidden, generated.ErrorResponse{Message: "unauthorized"})
	}

	token, err := Create(map[string]string{
		"sub": users.Slug,
	})
	if nil != err {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: http.StatusText(http.StatusInternalServerError)})
	}

	return ctx.JSON(
		http.StatusOK,
		generated.LoginResponse{
			Id:    users.Id,
			Token: token,
		},
	)
}

func (s *Server) Profile(ctx echo.Context) error {
	f := ctx.Get("user").(map[string]any)
	out, err := s.Repository.FindBySlug(ctx.Request().Context(), repository.FindBySlugInput{Slug: f["sub"].(string)})
	if nil != err {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: http.StatusText(http.StatusInternalServerError)})
	}

	return ctx.JSON(http.StatusOK, generated.ProfileResponse{
		FullName: out.FullName,
		Phone:    out.Phone,
	})
}

func (s *Server) UpdateProfile(ctx echo.Context) error {
	var (
		slug    = ctx.Get("user").(map[string]any)["sub"].(string)
		request generated.UpdateRequest
		c       = ctx.Request().Context()
	)

	existUser, err := s.Repository.FindByPhone(c, repository.FindByPhoneInput{Phone: request.Phone})
	if nil != err {
		if err != sql.ErrNoRows {
			return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: http.StatusText(http.StatusInternalServerError)})
		}
	}
	if existUser.Phone != "" {
		return ctx.JSON(http.StatusConflict, generated.ErrorResponse{Message: "phone number already exists"})
	}

	if err = s.Repository.Put(c, repository.UpdateUserInput{
		Slug:     slug,
		FullName: request.FullName,
		Phone:    request.Phone,
	}); nil != err {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: http.StatusText(http.StatusInternalServerError)})
	}

	return ctx.NoContent(http.StatusOK)
}

func (s *Server) Register(ctx echo.Context) error {
	var request generated.RegistrationRequest
	if err := ctx.Bind(&request); nil != err {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: "invalid parameter request"})
	}

	if err := validatePhone(request.Phone); nil != err {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	if err := validatePassword(request.Password); nil != err {
		return ctx.JSON(http.StatusBadRequest, generated.ErrorResponse{Message: err.Error()})
	}

	var c = ctx.Request().Context()
	users, err := s.Repository.FindByPhone(c, repository.FindByPhoneInput{Phone: request.Phone})
	if nil == err && !reflect.ValueOf(users).IsZero() {
		return ctx.JSON(http.StatusConflict, generated.ErrorResponse{Message: "user with phone number already exists"})
	}

	slug := base64.RawStdEncoding.EncodeToString([]byte(request.Phone))
	p, _ := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	out, err := s.Repository.Store(c, repository.RegistrationInput{
		Slug:     slug,
		FullName: request.FullName,
		Phone:    request.Phone,
		Password: string(p),
	})
	if nil != err {
		return ctx.JSON(http.StatusInternalServerError, generated.ErrorResponse{Message: http.StatusText(http.StatusInternalServerError)})
	}

	return ctx.JSON(http.StatusOK, generated.RegistrationResponse{Id: out.Id})
}

func validatePassword(password string) error {
	if len(password) < 6 && len(password) > 64 {
		return fmt.Errorf("password must be greater than 6 and lower than 64")
	}

	var (
		hasCapital = false
		hasNumber  = false
		hasSpecial = false
	)

	for _, r := range password {
		if unicode.IsUpper(r) {
			hasCapital = true
		}
		if unicode.IsNumber(r) {
			hasNumber = true
		}
		if !unicode.IsLetter(r) && !unicode.IsNumber(r) {
			hasSpecial = true
		}

		if hasCapital && hasNumber && hasSpecial {
			break
		}
	}

	if !hasCapital ||
		!hasNumber ||
		!hasSpecial {
		return fmt.Errorf("password must had minimum 1 capital letter, 1 number and 1 special character")
	}

	return nil
}

func validatePhone(phone string) error {
	if !strings.HasPrefix(phone, "+62") {
		return fmt.Errorf("phone number must have prefix +62")
	}
	phoneLen := len(strings.TrimPrefix(phone, "+62"))
	if phoneLen < 10 || phoneLen > 13 {
		return fmt.Errorf("phone number len must be greater than 10 and lower than 13")
	}

	return nil
}
