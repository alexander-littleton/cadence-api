package userservice

import (
	"context"
	"github.com/alexander-littleton/cadence-api/internal/models"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen --source=interfaces.go --destination=mocks/mock_dependencies.go --package=mocks
type UserRepository interface {
	CreateUser(ctx context.Context, user models.User) error
	GetUserById(ctx context.Context, userId primitive.ObjectID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type Validator interface {
	Struct(s interface{}) error
}
