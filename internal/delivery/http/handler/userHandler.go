package handler

import (
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type UserHandler struct {
	Log     *logrus.Logger
	UseCase *usecase.UserUseCase
}

func NewUserHandler(useCase *usecase.UserUseCase, logger *logrus.Logger) *UserHandler {
	return &UserHandler{
		Log:     logger,
		UseCase: useCase,
	}
}

func (c *UserHandler) Register(ctx *fiber.Ctx) error {
	request := new(model.RegisterUserRequest)
	err := ctx.BodyParser(request)
	if err != nil {
		c.Log.Warnf("Failed to parse request body : %+v", err)
		return fiber.ErrBadRequest
	}

	response, err := c.UseCase.Create(ctx.Context(), request)
	if err != nil {
		c.Log.Warnf("Failed to register user : %+v", err)
		return err
	}

	return ctx.JSON(model.WebResponse[*model.UserResponse]{Data: response})
}
