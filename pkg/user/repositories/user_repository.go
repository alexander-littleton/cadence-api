package repositories

import (
	"context"
	"github.com/alexander-littleton/cadence-api/pkg/user/domain"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		collection: collection,
	}
}

func (r *UserRepository) CreateUser(ctx context.Context, user domain.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) GetUserById(ctx context.Context, userId primitive.ObjectID) (domain.User, error) {
	user := &domain.User{}
	err := r.collection.FindOne(ctx, bson.D{{Key: "_id", Value: userId}}).Decode(user)
	if err != nil {
		return domain.User{}, err
	}
	//TODO: ensure empty user is handled during error catch
	return *user, nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (domain.User, error) {
	user := &domain.User{}
	err := r.collection.FindOne(ctx, bson.D{{Key: "email", Value: email}}).Decode(user)
	if err != nil {
		return domain.User{}, err
	}
	return *user, nil
}
