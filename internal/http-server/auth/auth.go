package auth

import (
	"UrlShorterService/internal/entity"
	"UrlShorterService/internal/http-server/response"
	"UrlShorterService/internal/services/jwt"
	"UrlShorterService/internal/services/user_service"
	"context"
	"github.com/go-chi/render"
	"log/slog"
	"net/http"
)

type registerRequest struct {
	Name     string `json:"name" validate:"required,alphanum"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum"`
}
type loginRequest struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,alphanum"`
}
type Response struct {
	message string
}
type UserRepositoryInterface interface {
	AddUser(ctx context.Context, user *entity.User) (int64, error)
	FindUserByEmail(ctx context.Context, email string) (*entity.User, error)
}

func Register(log *slog.Logger, repositoryInterface UserRepositoryInterface) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "auth.Register"
		var request registerRequest
		if err := render.DecodeJSON(r.Body, &request); err != nil {
			log.Error(op, "failed to decode body", err)
			render.JSON(w, r, response.Error("failed to decode body"))
			return
		}
		user, err := user_service.CreateUser(request.Name, request.Email, request.Password)
		if err != nil {
			log.Error(op, "failed to add user", err)
			render.JSON(w, r, response.Error("failed to register user"))
			return
		}
		id, err := repositoryInterface.AddUser(r.Context(), user)
		if err != nil {
			log.Error(op, "failed to add user", err)
			render.JSON(w, r, response.Error("failed to register user"))
			return
		}
		token, err := jwt.CreateJwt(id)
		render.JSON(w, r, response.OK(token))
	}
}
func Login(log *slog.Logger, repositoryInterface UserRepositoryInterface) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		const op = "auth.Login"
		var request loginRequest
		if err := render.DecodeJSON(r.Body, &request); err != nil {
			log.Error(op, "failed to decode body", err)
			render.JSON(w, r, "failed to decode body")
			return
		}
		user, err := repositoryInterface.FindUserByEmail(r.Context(), request.Email)
		if err != nil {
			render.JSON(w, r, response.Error("failed to find user"))
			return
		}
		if !user_service.IsCorrectPassword(request.Password, user.Salt, user.Password) {
			render.JSON(w, r, response.Error("Wrong password"))
			return
		}
		token, err := jwt.CreateJwt(user.Id)
		resp := response.OK(token)
		log.Info(op, resp)
		render.JSON(w, r, resp)
	}
}
