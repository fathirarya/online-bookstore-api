package usecase

import (
	"context"
	"time"

	"github.com/fathirarya/online-bookstore-api/internal/auth"
	"github.com/fathirarya/online-bookstore-api/internal/entity"
	"github.com/fathirarya/online-bookstore-api/internal/model"
	"github.com/fathirarya/online-bookstore-api/internal/model/converter"
	"github.com/fathirarya/online-bookstore-api/internal/repository"
	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserUseCase struct {
	DB             *gorm.DB
	Log            *logrus.Logger
	Validate       *validator.Validate
	UserRepository *repository.UserRepository
}

func NewUserUseCase(db *gorm.DB, logger *logrus.Logger, validate *validator.Validate,
	userRepository *repository.UserRepository) *UserUseCase {
	return &UserUseCase{
		DB:             db,
		Log:            logger,
		Validate:       validate,
		UserRepository: userRepository,
	}
}

// Register registers a new user with transaction
func (uc *UserUseCase) Register(ctx context.Context, req *model.RegisterUserRequest) (*model.UserResponse, error) {
	if err := uc.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	// Cek email sudah terdaftar
	existingUser, err := uc.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil && err != gorm.ErrRecordNotFound {
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}
	if existingUser != nil {
		return nil, fiber.NewError(fiber.StatusConflict, "email already registered")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		tx.Rollback()
		uc.Log.Error("failed to hash password: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	user := &entity.User{
		Name:     req.Name,
		Email:    req.Email,
		Password: string(hashedPassword),
	}

	if err := uc.UserRepository.Create(tx, user); err != nil {
		tx.Rollback()
		uc.Log.Error("failed to create user: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to register user")
	}

	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to register user")
	}

	return converter.UserToResponse(user), nil
}

// Login authenticates a user and returns AuthResponse with JWT token
func (uc *UserUseCase) Login(ctx context.Context, req *model.LoginUserRequest, jwtService *auth.JWTService) (*model.AuthResponse, error) {
	if err := uc.Validate.Struct(req); err != nil {
		return nil, fiber.NewError(fiber.StatusBadRequest, err.Error())
	}

	tx := uc.DB.WithContext(ctx).Begin()
	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
		}
	}()

	user, err := uc.UserRepository.FindByEmail(ctx, req.Email)
	if err != nil {
		tx.Rollback()
		if err == gorm.ErrRecordNotFound {
			return nil, fiber.NewError(fiber.StatusBadRequest, "invalid email or password")
		}
		uc.Log.Error("failed to find user: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "internal server error")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(req.Password)); err != nil {
		tx.Rollback()
		return nil, fiber.NewError(fiber.StatusBadRequest, "invalid email or password")
	}

	token, err := jwtService.GenerateToken(user.ID)
	if err != nil {
		tx.Rollback()
		uc.Log.Error("failed to generate token: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to generate token")
	}

	expiresAt := time.Now().Add(jwtService.ExpireDuration())

	if err := tx.Commit().Error; err != nil {
		uc.Log.Error("failed to commit transaction: ", err)
		return nil, fiber.NewError(fiber.StatusInternalServerError, "failed to login user")
	}

	return converter.AuthToResponse(user, token, expiresAt), nil
}
