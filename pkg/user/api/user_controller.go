package api

import (
	"errors"
	"fmt"
	"github.com/alexander-littleton/cadence-api/pkg/common/cadence_errors"
	userService "github.com/alexander-littleton/cadence-api/pkg/user"
	"github.com/alexander-littleton/cadence-api/pkg/user/domain"
	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
)

type Controller struct {
	userService userService.Service
}

func New(userService userService.Service) Controller {
	return Controller{
		userService: userService,
	}
}

func (r Controller) RegisterRoutes(router *gin.Engine) {
	router.POST("/user", r.createUser)
	router.GET("/user/:email", r.GetUserByEmail)
	//router.PUT("/user/:userId", userController.EditAUser())
	//router.DELETE("/user/:userId", userController.DeleteAUser())
	//router.GET("/users", controllers.GetAllUsers())
}

func (r Controller) createUser(ctx *gin.Context) {
	var newUser domain.User
	fmt.Printf("%+v\n", ctx)
	if err := ctx.BindJSON(&newUser); err != nil {
		ctx.JSON(http.StatusBadRequest, domain.UserResponse{
			Status:  http.StatusBadRequest,
			Message: "error",
			Data: map[string]interface{}{
				"data": fmt.Sprint("failed to unmarshal new user from request body: ", err.Error()),
			},
		})
		return
	}

	createdUser, err := r.userService.CreateUser(ctx, newUser)
	if err != nil {
		var status int
		if errors.Is(err, cadence_errors.ValidationErr) {
			status = http.StatusBadRequest
		} else {
			status = http.StatusInternalServerError
		}
		ctx.JSON(
			status,
			domain.UserResponse{
				Status:  status,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			},
		)
		return
	}

	ctx.JSON(
		http.StatusCreated,
		domain.UserResponse{
			Status:  http.StatusCreated,
			Message: "success",
			Data:    map[string]interface{}{"data": createdUser},
		},
	)
	return
}

func (r Controller) GetUserById(ctx *gin.Context) {
	rawId := ctx.Param("userId")
	objId, _ := primitive.ObjectIDFromHex(rawId)

	user, err := r.userService.GetUserById(ctx, objId)
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			domain.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		domain.UserResponse{
			Status:  http.StatusOK,
			Message: "success",
			Data:    map[string]interface{}{"data": user},
		},
	)
}

func (r Controller) GetUserByEmail(ctx *gin.Context) {
	email := ctx.Param("email")

	user, err := r.userService.GetUserByEmail(ctx, email)
	//TODO: handle validation errors
	if err != nil {
		ctx.JSON(
			http.StatusInternalServerError,
			domain.UserResponse{
				Status:  http.StatusInternalServerError,
				Message: "error",
				Data:    map[string]interface{}{"data": err.Error()},
			},
		)
		return
	}

	ctx.JSON(
		http.StatusOK,
		domain.UserResponse{
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
