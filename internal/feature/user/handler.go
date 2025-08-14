package user

import (
	"context"
	"errors"

	"log/slog"
	"net/http"

	"github.com/Sanchir01/go-shortener/pkg/logger"
	"github.com/Sanchir01/go-shortener/pkg/utils"
	"github.com/go-playground/validator/v10"
	"github.com/google/uuid"

	"github.com/Sanchir01/go-shortener/pkg/api"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
)

type Handler struct {
	Service HandlerUser
	log     *slog.Logger
}

//go:generate go run github.com/vektra/mockery/v2@v2.52.2 --name=HandlerUser
type HandlerUser interface {
	Register(ctx context.Context, p RegisterParams) (*uuid.UUID, error)
	Login(ctx context.Context, email, password string) (*DatabaseUser, error)
}

func NewHandler(s HandlerUser, lg *slog.Logger) *Handler {
	return &Handler{
		Service: s,
		log:     lg,
	}
}

// @Summary Auth
// @Tags auth
// @Description register user
// @Accept json
// @Produce json
// @Param input body LoginRequest true "login body"
// @Success 201 {object}  AuthResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /register [post]
func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	const op = "User.Handler.Register"
	log := h.log.With(
		slog.String("op", op),
		slog.String("request_id", middleware.GetReqID(r.Context())),
	)
	var req AuthRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	id, err := h.Service.Register(r.Context(), RegisterParams{
		Email:    &req.Email,
		Title:    req.Username,
		Password: &req.Password,
	})
	if errors.Is(err, utils.ErrorUserAlreadyExists) {
		log.Error("user already exists", logger.Err(err))
		render.Status(r, http.StatusBadRequest)
		render.JSON(w, r, api.Error("Пользователь с таким email или username уже существует"))
		return
	}
	if err != nil {
		log.Error("failed to register user", logger.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("ошибка регистрации"))
		return
	}
	log.Info("login success")

	if err = AddCookieTokens(*id, w, "user", "localhost"); err != nil {
		log.Error("register cookie errors", logger.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to register cookie"))
		return
	}
	render.Status(r, http.StatusCreated)
	render.JSON(w, r, AuthResponse{
		Response: api.OK(),
		Email:    req.Email,
		Username: req.Username,
	})
}

// @Summary Login
// @Tags auth
// @Description login user
// @Accept json
// @Produce json
// @Param input body LoginRequest true "auth body"
// @Success 200 {object}  LoginResponse
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /login [post]
func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	const op = "User.Handler.Login"
	log := h.log.With(slog.String("op", op))
	var req LoginRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", logger.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}
	user, err := h.Service.Login(r.Context(), req.Email, req.Password)
	if errors.Is(err, utils.ErrorUserNotFound) {
		log.Error("invalid login", logger.Err(err))
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, api.Error("Неправильный логин"))
		return
	}
	if errors.Is(err, utils.ErrorInvalidPassword) {
		log.Error("invalid password", logger.Err(err))
		render.Status(r, http.StatusNotFound)
		render.JSON(w, r, api.Error("Неправильный логин"))
		return
	}
	if err != nil {
		log.Error("failed to login", logger.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("internal server error"))
		return
	}
	if err = AddCookieTokens(user.ID, w, user.Role, "localhost"); err != nil {
		log.Error("register cookie errors", logger.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to register cookie"))
		return
	}
	render.JSON(w, r, LoginResponse{
		Response: api.OK(),
		Email:    user.Email,
		Username: user.Name,
	})
}

// @Summary GogleLogin
// @Tags auth
// @Description google login user
// @Accept json
// @Produce json
// @Success 200 {object}  api.Response
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /google [get]
func (h *Handler) GoogleLogin(w http.ResponseWriter, r *http.Request) {
	const op = "User.Handler.GoogleLogin"
	log := h.log.With(slog.String("op", op))

	url, err := utils.GetUrlGoogleString()
	if err != nil {
		log.Error("failed to generate google oauth url", logger.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("internal server error"))
		return
	}

	log.Info("success get url", slog.String("url", url))
	http.Redirect(w, r, url, http.StatusTemporaryRedirect)
}

// @Summary GogleRegister
// @Tags auth
// @Description google register user
// @Accept json
// @Produce json
// @Param input body GoogleRegisterRequest true "auth body"
// @Success 200 {object}  api.Response
// @Failure 400,404 {object}  api.Response
// @Failure 500 {object}  api.Response
// @Router /google [get]
func (h *Handler) GoogleCallback(w http.ResponseWriter, r *http.Request) {
	const op = "User.Handler.GoogleCallback"
	log := h.log.With(slog.String("op", op))

	var req GoogleRegisterRequest
	if err := render.DecodeJSON(r.Body, &req); err != nil {
		log.Error("failed to decode request body", slog.Any("err", err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("Ошибка при валидации данных"))
		return
	}
	if err := validator.New().Struct(req); err != nil {
		log.Error("invalid request", logger.Err(err))
		render.JSON(w, r, api.Error("invalid request"))
		return
	}

	// Exchange code for token
	data, err := utils.ExchangeGoogleCodeForToken(r.Context(), req.Code)
	if err != nil {
		log.Error("failed to exchange code for token", logger.Err(err))
		render.Status(r, http.StatusBadGateway)
		render.JSON(w, r, api.Error("failed to exchange authorization code"))
		return
	}

	log.Info("successfully exchanged code for token")
	if _, err := h.Service.Register(r.Context(), RegisterParams{
		Email: &data.Email,
		Title: data.Name,
	}); err != nil {
		log.Error("failed to register user", logger.Err(err))
		render.Status(r, http.StatusInternalServerError)
		render.JSON(w, r, api.Error("failed to register user"))
		return
	}

	render.Status(r, http.StatusCreated)
	render.JSON(w, r, map[string]interface{}{
		"status":  "success",
		"message": "Google authentication successful",
	})
}
