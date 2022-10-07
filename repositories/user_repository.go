package repositories

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"internal/configs"
	"internal/models"
)

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository() *UserRepository {
	return &UserRepository{
		collection: configs.GetCollection(configs.DB, "users"),
	}
}

func (r *UserRepository) Create(ctx context.Context, user models.User) error {
	_, err := r.collection.InsertOne(ctx, user)
	if err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) Find(ctx context.Context, userId primitive.ObjectID) error {
	user := &models.User{}
	err := r.collection.FindOne(ctx, bson.D{{"_id", userId}}).Decode(user)
	if err != nil {
		return err
	}
	return nil
}
