package userService

import (
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"internal/models"
	"internal/repositories"
)

var validate = validator.New()

type UserService struct {
	userRepository *repositories.UserRepository
}

func NewUserService(userRepo *repositories.UserRepository) *UserService {
	return &UserService{
		userRepository: userRepo,
	}
}

func (r *UserService) CreateUser(ctx *gin.Context, user models.User) (models.User, error) {
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

func (r *UserService) GetUser(ctx *gin.Context, userId primitive.ObjectID) (*models.User, error) {
	user, err := r.userRepository.Find(ctx, userId)
	if err != nil {
		return nil, fmt.Errorf("failed to get user with id %s: %w", userId.Hex(), err)
	}
	return user, nil
}
