package domain

import "go.mongodb.org/mongo-driver/bson/primitive"

type User struct {
	Id    primitive.ObjectID `json:"id" bson:"_id"`
	Email string             `json:"email,omitempty" validate:"required"`
}

type UserResponse struct {
	Status  int                    `json:"status"`
	Message string                 `json:"message"`
	Data    map[string]interface{} `json:"data"`
}

type CreateUserRequest struct {
	Email string `json:"email,omitempty" validate:"required"`
}
