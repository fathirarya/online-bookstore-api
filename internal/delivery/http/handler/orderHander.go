package handler

import (
	"strconv"

	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/usecase"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
)

type OrderHandler struct {
	Log     *logrus.Logger
	UseCase *usecase.OrderUseCase
}

func NewOrderHandler(useCase *usecase.OrderUseCase, logger *logrus.Logger) *OrderHandler {
	return &OrderHandler{
		Log:     logger,
		UseCase: useCase,
	}
}

func (h *OrderHandler) Create(ctx *fiber.Ctx) error {
	var request model.CreateOrderRequest
	if err := ctx.BodyParser(&request); err != nil {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "Invalid request body",
		})
	}

	// Ambil userID dari token (middleware)
	userID, ok := ctx.Locals("user_id").(int)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ValidationErrorResponse{
			Message: "Unauthorized, user not found",
		})
	}

	// Panggil UseCase untuk membuat order
	response, err := h.UseCase.CreateOrder(ctx.Context(), &request, userID)
	if err != nil {
		// Handle fiber.Error dari UseCase
		if fiberErr, ok := err.(*fiber.Error); ok {
			// BadRequest
			if fiberErr.Code == fiber.StatusBadRequest {
				msg := fiberErr.Message
				errorsMap := make(map[string]string)

				if msg == "validation failed, please check your input" {
					errorsMap["body"] = "Please check your input fields"
					msg = "Request validation failed"
				} else if msg == "maximum 5 books per transaction" {
					errorsMap["items"] = "You can order maximum 5 books per transaction"
					msg = "Validation failed"
				} else if msg == "book not found" {
					errorsMap["book"] = "One or more books not found"
					msg = "Validation failed"
				}

				return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
					Message: msg,
					Errors:  errorsMap,
				})
			}

			// InternalServerError
			if fiberErr.Code == fiber.StatusInternalServerError {
				return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
					Message: "Internal server error, please try again later",
				})
			}

			// Default fallback
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}

		// Non-fiber error fallback
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "Internal server error",
		})
	}

	// Response sukses
	return ctx.Status(fiber.StatusCreated).JSON(model.WebResponse[*model.OrderResponse]{
		Data: response,
	})
}

func (h *OrderHandler) Pay(ctx *fiber.Ctx) error {
	// 1. Ambil orderID dari URL param
	orderIDStr := ctx.Params("id")
	orderID, err := strconv.Atoi(orderIDStr)
	if err != nil || orderID <= 0 {
		return ctx.Status(fiber.StatusBadRequest).JSON(model.ValidationErrorResponse{
			Message: "Invalid order ID",
		})
	}

	// 2. Ambil userID dari token (middleware)
	userID, ok := ctx.Locals("user_id").(int)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ValidationErrorResponse{
			Message: "Unauthorized, user not found",
		})
	}

	// 3. Panggil UseCase untuk membayar order
	response, err := h.UseCase.PayOrder(ctx.Context(), orderID, userID)
	if err != nil {
		// Handle fiber.Error dari UseCase
		if fiberErr, ok := err.(*fiber.Error); ok {
			errorsMap := make(map[string]string)
			msg := fiberErr.Message

			if fiberErr.Code == fiber.StatusBadRequest {
				if msg == "order is not pending" {
					errorsMap["status"] = "Only pending orders can be paid"
					msg = "Validation failed"
				}
				return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
					Message: msg,
					Errors:  errorsMap,
				})
			}

			if fiberErr.Code == fiber.StatusForbidden {
				return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
					Message: "You are not allowed to pay this order",
				})
			}

			if fiberErr.Code == fiber.StatusNotFound {
				return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
					Message: "Order not found",
				})
			}

			if fiberErr.Code == fiber.StatusInternalServerError {
				return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
					Message: "Internal server error, please try again later",
				})
			}

			// Default fallback
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}

		// Non-fiber error fallback
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "Internal server error",
		})
	}

	// 4. Response sukses
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.OrderResponse]{
		Data: response,
	})
}

func (h *OrderHandler) List(ctx *fiber.Ctx) error {
	// Ambil userID dari token (middleware)
	userID, ok := ctx.Locals("user_id").(int)
	if !ok {
		return ctx.Status(fiber.StatusUnauthorized).JSON(model.ValidationErrorResponse{
			Message: "Unauthorized, user not found",
		})
	}

	// Panggil UseCase untuk ambil daftar order
	response, err := h.UseCase.GetOrdersByUser(ctx.Context(), userID)
	if err != nil {
		if fiberErr, ok := err.(*fiber.Error); ok {
			// InternalServerError
			if fiberErr.Code == fiber.StatusInternalServerError {
				return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
					Message: "Internal server error, please try again later",
				})
			}

			// Default fallback
			return ctx.Status(fiberErr.Code).JSON(model.ValidationErrorResponse{
				Message: fiberErr.Message,
			})
		}

		// Non-fiber error fallback
		return ctx.Status(fiber.StatusInternalServerError).JSON(model.ValidationErrorResponse{
			Message: "Internal server error",
		})
	}

	// Response sukses
	return ctx.Status(fiber.StatusOK).JSON(model.WebResponse[*model.OrderListResponse]{
		Data: response,
	})
}
