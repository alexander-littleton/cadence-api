package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"internal/models"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(collection *mongo.Collection) *UserRepository {
	return &UserRepository{
		collection: collection,
	}
}

func (r *UserRepository) Create(ctx context.Context, user models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Find(ctx context.Context, userId primitive.ObjectID) (*models.User, error) {
	user := &models.User{}
	err := r.collection.FindOne(ctx, bson.D{{"_id", userId}}).Decode(user)
	if err != nil {
		return nil, err
	}
	return user, nil
}
