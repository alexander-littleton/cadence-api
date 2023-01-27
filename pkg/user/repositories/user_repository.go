package repositories

import (
	"context"
	"github.com/alexander-littleton/cadence-api/pkg/user/domain"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

//go:generate mockgen --source=user_repository.go --destination=mocks/mock_dependencies.go --package=mocks
type UserRepository interface {
	CreateUser(ctx context.Context, user domain.User) error
	GetUserById(ctx context.Context, userId primitive.ObjectID) (domain.User, error)
	GetUserByEmail(ctx context.Context, email string) (domain.User, error)
}
