package userService

import (
	"context"
	"errors"
	"fmt"

	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"internal/models"
)

var validate = validator.New()

//go:generate mockgen --source=user_service.go --destination=mocks/mock_user_repository.go --package=mocks UserRepository
type UserRepository interface {
	Create(ctx context.Context, user models.User) error
	Find(ctx context.Context, userId primitive.ObjectID) (*models.User, error)
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
	//TODO: ensure we can call this if Id is blank?
	if !user.Id.IsZero() {
		return models.User{}, errors.New("expected a user without an id")
	}

	user.Id = primitive.NewObjectID()

	err := validate.Struct(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("%s: %w", NewUserValidationErr, err)
	}

	err = r.userRepository.Create(ctx, user)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return user, nil
}

func (r *UserService) GetUser(ctx context.Context, userId primitive.ObjectID) (*models.User, error) {
	user, err := r.userRepository.Find(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with id %s: %w", userId.Hex(), err)
	}
	return user, nil
}
