package controllers

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"strings"
	"time"

	"internal/models"
	"internal/responses"
	userService "internal/services"
)

var validate = validator.New()

//go:generate mockgen --source=controllers.go --destination=mocks/mock_user_service.go --package=mocks UserService
type UserService interface {
	CreateUser(ctx context.Context, user models.User) (models.User, error)
	GetUser(ctx context.Context, userId primitive.ObjectID) (*models.User, error)
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
	var user CreateUserRequest
	//TODO: validate the request body
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

	validate := validator.New()
	err := validate.Struct(user)
	if err != nil {
		fmt.Println(err)
		ctx.JSON(http.StatusBadRequest, responses.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: map[string]interface{}{
				"data": fmt.Sprint("failed to validate user: ", err.Error()),
			},
		})
	}
	createdUser, err := r.userService.CreateUser(ctx, user)
	if err != nil {
		var status int
		//TODO: try to switch this over to a errors.in()
		if strings.Contains(err.Error(), userService.NewUserValidationErr) {
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

func (r *UserController) GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		rawId := c.Param("userId")
		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(rawId)

		user, err := r.userService.GetUser(ctx, objId)
		if err != nil {
			c.JSON(
				http.StatusInternalServerError,
				responses.UserResponse{
					Status:  http.StatusInternalServerError,
					Message: "error",
					Data:    map[string]interface{}{"data": err.Error()},
				},
			)
			return
		}

		c.JSON(
			http.StatusOK,
			responses.UserResponse{
				Status:  http.StatusOK,
				Message: "success",
				Data:    map[string]interface{}{"data": user},
			},
		)
	}
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
//		results, err := userCollection.Find(ctx, bson.M{})
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
