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
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid request body",
		})
	}

	response, err := h.UseCase.Register(ctx.Context(), &request)
	if err != nil {
		// Tangkap error dari usecase
		if fiberErr, ok := err.(*fiber.Error); ok {
			// Cek jika error "email already registered"
			if fiberErr.Code == fiber.StatusConflict && fiberErr.Message == "email already registered" {
				return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
					Message: "validation failed",
					Errors:  map[string]string{"email": "email already registered"},
				})
			}
			// Error lain, misal internal server error
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}
		// Error tidak diketahui
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{
		Data: response,
	})
}

func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var request model.LoginUserRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid request body",
		})
	}

	response, err := h.UseCase.Login(ctx.Context(), &request, h.JWTService)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			// Mapping error validasi login
			if fiberErr.Code == fiber.StatusBadRequest && fiberErr.Message == "invalid email or password" {
				return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
					Message: "validation failed",
					Errors:  map[string]string{"email": "invalid email or password", "password": "invalid email or password"},
				})
			}
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	return ctx.JSON(model.WebResponse[*model.AuthResponse]{
		Data: response,
	})
}
