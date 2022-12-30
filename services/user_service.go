package userService

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"internal/models"
)

var validate = validator.New()

//go:generate mockgen --source=user_service.go --destination=mocks/mock_user_repository.go --package=mocks UserRepository
type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserById(ctx context.Context, userId primitive.ObjectID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type UserService struct {
	userRepository UserRepository
}

func NewUserService(userRepo UserRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

func (r *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	validatedUser, err := r.validateNewUser(ctx, user)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", NewUserValidationErr, err)
	}

	err = r.userRepository.CreateUser(ctx, validatedUser)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return validatedUser, nil
}

func (r *UserService) validateNewUser(ctx context.Context, user models.User) (models.User, error) {
	if !user.Id.IsZero() {
		return models.User{}, errors.New("expected a user without an id")
	}

	user.Id = primitive.NewObjectID()

	if _, err := mail.ParseAddress(user.Email); err != nil {
		return models.User{}, errors.New("invalid email address")
	}

	_, err := r.GetUserByEmail(ctx, user.Email)
	if err != nil {
	}

	err = validate.Struct(&user)
	if err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (r *UserService) GetUserById(ctx context.Context, userId primitive.ObjectID) (models.User, error) {
	if userId.IsZero() {
		return models.User{}, errors.New("valid user id must be provided")
	}
	user, err := r.userRepository.GetUserById(ctx, userId)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user with id %s: %w", userId.Hex(), err)
	}
	return user, nil
}

func (r *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	user, err := r.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user with email %s: %w", email, err)
	}
	return user, nil
}
