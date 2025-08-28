package handler

import (
	"github.com/fathirarya/online-bookstore-api/internal/auth"
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/usecase"
	"github.com/fathirarya/online-bookstore-api/internal/utils"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Log        *logrus.Logger
	UseCase    *usecase.UserUseCase
	JWTService *auth.JWTService
	Validate   *validator.Validate
}

func NewUserHandler(useCase *usecase.UserUseCase, logger *logrus.Logger, jwtService *auth.JWTService, validate *validator.Validate) *UserHandler {
	return &UserHandler{
		Log:        logger,
		UseCase:    useCase,
		JWTService: jwtService,
		Validate:   validate,
	}
}

func (h *UserHandler) Register(ctx *fiber.Ctx) error {
	var request model.RegisterUserRequest

	// Parse request body into RegisterUserRequest struct
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid request body",
		})
	}

	// Validate input using struct tags (e.g., `required`, `email`)
	if err := h.Validate.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "validation failed",
			Errors:  utils.TranslateValidationErrors(err), // convert validator errors to a field:error map
		})
	}

	// Call UseCase to handle user registration
	response, err := h.UseCase.Register(ctx.Context(), &request)
	if err != nil {
		// Handle error returned from UseCase
		if fiberErr, ok := err.(*fiber.Error); ok {
			// Case: email already registered
			if fiberErr.Code == fiber.StatusConflict && fiberErr.Message == "email already registered" {
				return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
					Message: "validation failed",
					Errors:  map[string]string{"email": "email already registered"},
				})
			}

			// Other known errors, return directly
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}

		// Unknown error, return generic internal server error
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	// Return success response with created user data
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.UserResponse]{
		Data: response,
	})
}

// Login handles user authentication and JWT token generation
func (h *UserHandler) Login(ctx *fiber.Ctx) error {
	var request model.LoginUserRequest

	// 1Parse request body into LoginUserRequest struct
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "invalid request body",
		})
	}

	// Validate input using struct tags (e.g., `required`, `email`)
	if err := h.Validate.Struct(request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "validation failed",
			Errors:  utils.TranslateValidationErrors(err), // helper translate to map field:error
		})
	}

	// Call UseCase to perform login and generate JWT
	response, err := h.UseCase.Login(ctx.Context(), &request, h.JWTService)
	if err != nil {
		// Handle error returned from UseCase
		if fiberErr, ok := err.(*fiber.Error); ok {
			// Case: invalid credentials (email or password)
			if fiberErr.Code == fiber.StatusBadRequest && fiberErr.Message == "invalid email or password" {
				return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
					Message: "validation failed",
					Errors: map[string]string{
						"email":    "invalid email or password",
						"password": "invalid email or password",
					},
				})
			}

			// Other known errors, return directly
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}

		// Unknown error, return generic internal server error
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "internal server error",
		})
	}

	// Return success response with JWT token and user info
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.AuthResponse]{
		Data: response,
	})
}
