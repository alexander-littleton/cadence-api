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
}

func (r Controller) createUser(ctx *gin.Context) {
	var newUser domain.User
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
