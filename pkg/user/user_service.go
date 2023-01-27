package user

import (
	"context"
	"errors"
	"fmt"
	"net/mail"

	"github.com/alexander-littleton/cadence-api/pkg/common/cadence_errors"
	"github.com/alexander-littleton/cadence-api/pkg/user/domain"
	"github.com/alexander-littleton/cadence-api/pkg/user/repositories"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen --source=user_service.go --destination=mocks/mock_user_service.go --package=mocks
type Service interface {
	CreateUser(ctx context.Context, user domain.User) (domain.User, error)
	GetUserById(ctx context.Context, userId primitive.ObjectID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}

type service struct {
	userRepository repositories.UserRepository
}

func New(userRepo repositories.UserRepository) Service {
	return &service{
		userRepository: userRepo,
	}
}

func (r *service) CreateUser(ctx context.Context, user domain.User) (domain.User, error) {
	validatedUser, err := r.validateNewUser(ctx, user)
	if err != nil {
		return domain.User{}, err
	}

	validatedUser.Id = primitive.NewObjectID()

	err = r.userRepository.CreateUser(ctx, validatedUser)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to create user: %w", err)
	}

	return validatedUser, nil
}

func (r *service) validateNewUser(ctx context.Context, user domain.User) (domain.User, error) {
	if !user.Id.IsZero() {
		return domain.User{}, fmt.Errorf("%w: %s", cadence_errors.ValidationErr, "expected a user without an id")
	}

	_, err := r.GetUserByEmail(ctx, user.Email)
	if err != nil && !errors.Is(err, cadence_errors.ErrNotFound) {
		return domain.User{}, fmt.Errorf("%s: %w", "failed to get user by email", err)
	} else if err == nil {
		return domain.User{}, fmt.Errorf("%w: %s", cadence_errors.ValidationErr, "user with email already exists")
	}

	validate := validator.Validate{}
	err = validate.Struct(&user)
	if err != nil {
		return domain.User{}, fmt.Errorf("%w: %s", cadence_errors.ValidationErr, err.Error())
	}
	return user, nil
}

func (r *service) GetUserById(ctx context.Context, userId primitive.ObjectID) (domain.User, error) {
	if userId.IsZero() {
		return domain.User{}, fmt.Errorf("%w: %s", cadence_errors.ValidationErr, "valid user id must be provided")
	}
	user, err := r.userRepository.GetUserById(ctx, userId)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user with id %s: %w", userId.Hex(), err)
	}
	return user, nil
}

//GetUserByEmail takes an email, validates it, then returns the user with matching email. If a user does not exist in
//the db, then it will return an error.
func (r *service) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	if _, err := mail.ParseAddress(email); err != nil {
		return domain.User{}, fmt.Errorf("%w: %s", cadence_errors.ValidationErr, err.Error())
	}

	user, err := r.userRepository.GetUserByEmail(ctx, email)
	if err != nil {
		return domain.User{}, fmt.Errorf("failed to get user with email %s: %w", email, err)
	}
	return user, nil
}
