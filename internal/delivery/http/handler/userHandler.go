package handler

import (
	"github.com/fathirarya/online-bookstore-api/internal/auth"
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Log        *logrus.Logger
	UseCase    *usecase.UserUseCase
	JWTService *auth.JWTService
}

func NewUserHandler(useCase *usecase.UserUseCase, logger *logrus.Logger, jwtService *auth.JWTService) *UserHandler {
	return &UserHandler{
		Log:        logger,
		UseCase:    useCase,
		JWTService: jwtService,
	}
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	var request model.RegisterUserRequest
	if err := ctx.BodyParser(&request); err != nil {
		h.Log.Warnf("Failed to parse request body: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := h.UseCase.Register(ctx.Context(), &request)
	if err != nil {
		h.Log.Warnf("Failed to register user: %+v", err)
		return err
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{
		Data: response,
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var request model.LoginUserRequest
	if err := ctx.BodyParser(&request); err != nil {
		h.Log.Warnf("Failed to parse login request: %+v", err)
		return fiber.NewError(fiber.StatusBadRequest, "invalid request body")
	}

	response, err := h.UseCase.Login(ctx.Context(), &request, h.JWTService)
	if err != nil {
		h.Log.Warnf("Failed to login user: %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.AuthResponse]{
		Data: response,
	})
}
