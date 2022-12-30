package userservice

import (
	"context"
	"errors"
	"fmt"
	. "github.com/alexander-littleton/cadence-api/internal/common/cadence_errors"
	"github.com/alexander-littleton/cadence-api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/mail"
)

type UserService struct {
	userRepository UserRepository
	validator      Validator
}

func NewUserService(userRepo UserRepository, validator Validator) *UserService {
	return &UserService{
		userRepository: userRepo,
		validator:      validator,
	}
}

func (r *UserService) CreateUser(ctx context.Context, user models.User) (models.User, error) {
	validatedUser, err := r.validateNewUser(ctx, user)
	if err != nil {
		return models.User{}, err
	}

	validatedUser.Id = primitive.NewObjectID()

	err = r.userRepository.CreateUser(ctx, validatedUser)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return validatedUser, nil
}

func (r *UserService) validateNewUser(ctx context.Context, user models.User) (models.User, error) {
	if !user.Id.IsZero() {
		return models.User{}, fmt.Errorf("%w: %s", ValidationErr, "expected a user without an id")
	}

	_, err := r.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, ErrNotFound) {
		return models.User{}, fmt.Errorf("%s: %w", "failed to get user by email", err)
	} else if err == nil {
		return models.User{}, fmt.Errorf("%w: %s", ValidationErr, "user with email already exists")
	}

	err = r.validator.Struct(&user)
	if err != nil {
		return models.User{}, fmt.Errorf("%w: %s", ValidationErr, err.Error())
	}
	return user, nil
}

func (r *UserService) GetUserById(ctx context.Context, userId primitive.ObjectID) (models.User, error) {
	if userId.IsZero() {
		return models.User{}, fmt.Errorf("%w: %s", ValidationErr, "valid user id must be provided")
	}
	user, err := r.userRepository.GetUserById(ctx, userId)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user with id %s: %w", userId.Hex(), err)
	}
	return user, nil
}

func (r *UserService) GetUserByEmail(ctx context.Context, email string) (models.User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return models.User{}, fmt.Errorf("%w: %s", ValidationErr, err.Error())
	}

	user, err := r.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return models.User{}, fmt.Errorf("failed to get user with email %s: %w", email, err)
	}
	return user, nil
}
