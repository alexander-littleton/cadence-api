package controllers

import (
	"context"
	"errors"
	"fmt"
	"github.com/alexander-littleton/cadence-api/internal/common/cadence_errors"
	"github.com/alexander-littleton/cadence-api/internal/models"
	"github.com/alexander-littleton/cadence-api/internal/responses"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

//go:generate mockgen --source=controllers.go --destination=mocks/mock_user_service.go --package=mocks UserService
type UserService interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUserById(ctx context.Context, userId primitive.ObjectID) (models.User, error)
	GetUserByEmail(ctx context.Context, email string) (models.User, error)
}

type UserController struct {
	userService UserService
}

func NewUserController(userService UserService) *UserController {
	return &UserController{
		userService: userService,
	}
}

type CreateUserRequest struct {
	Email string `json:"email,omitempty" validate:"required"`
}

func (r *UserController) CreateUser(ctx *gin.Context) {
	var user models.User
	if err := ctx.BindJSON(&user); err != nil {
		ctx.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: map[string]interface{}{
				"data": fmt.Sprint("failed to unmarshal user from request body: ", err.Error()),
			},
		})
		return
	}

	createdUser, err := r.userService.CreateUser(ctx, user)
	if err != nil {
		var status int
		if errors.Is(err, cadence_errors.ValidationErr) {
			status = http.StatusBadRequest
		} else {
			status = http.StatusInternalServerError
		}
		ctx.JSON(
			status,
			responses.UserResponse{
				Status:  status,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			},
		)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		responses.UserResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data:    map[string]interface{}{"data": createdUser},
		},
	)
	return
}

func (r *UserController) GetUserById(ctx *gin.Context) {
	rawId := ctx.Param("userId")
	objId, _ := primitive.ObjectIDFromHex(rawId)

	user, err := r.userService.GetUserById(ctx, objId)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"data": user},
		},
	)
}

func (r *UserController) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	user, err := r.userService.GetUserByEmail(ctx, email)
	//TODO: handle validation errors
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			responses.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		responses.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"data": user},
		},
	)
}

//func EditAUser() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//		userId := c.Param("userId")
//		var user models.User
//		defer cancel()
//
//		objId, _ := primitive.ObjectIDFromHex(userId)
//
//		//validate the request body
//		if err := c.BindJSON(&user); err != nil {
//			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
//			return
//		}
//
//		//use the validator library to validate required fields
//		if validationErr := validate.Struct(&user); validationErr != nil {
//			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
//			return
//		}
//
//		update := bson.M{"email": user.Email}
//		result, err := userCollection.UpdateOne(ctx, bson.M{"id": objId}, bson.M{"$set": update})
//
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
//			return
//		}
//
//		//get updated user details
//		var updatedUser models.User
//		if result.MatchedCount == 1 {
//			err := userCollection.FindOne(ctx, bson.M{"id": objId}).Decode(&updatedUser)
//			if err != nil {
//				c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
//				return
//			}
//		}
//
//		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": updatedUser}})
//	}
//}
//
//func DeleteAUser() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//		userId := c.Param("userId")
//		defer cancel()
//
//		objId, _ := primitive.ObjectIDFromHex(userId)
//
//		result, err := userCollection.DeleteOne(ctx, bson.M{"id": objId})
//
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
//			return
//		}
//
//		if result.DeletedCount < 1 {
//			c.JSON(http.StatusNotFound,
//				responses.UserResponse{Status: http.StatusNotFound, Message: "error", Data: map[string]interface{}{"data": "User with specified ID not found!"}},
//			)
//			return
//		}
//
//		c.JSON(http.StatusOK,
//			responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": "User successfully deleted!"}},
//		)
//	}
//}

//func GetAllUsers() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
//		var users []models.User
//		defer cancel()
//
//		results, err := userCollection.GetUserById(ctx, bson.M{})
//
//		if err != nil {
//			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
//			return
//		}
//
//		//reading from the db in an optimal way
//		defer results.Close(ctx)
//		for results.Next(ctx) {
//			var singleUser models.User
//			if err = results.Decode(&singleUser); err != nil {
//				c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
//			}
//
//			users = append(users, singleUser)
//		}
//
//		c.JSON(http.StatusOK,
//			responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": users}},
//		)
//	}
//}
