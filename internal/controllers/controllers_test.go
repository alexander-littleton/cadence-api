package controllers_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alexander-littleton/cadence-api/internal/common/cadence_errors"
	"github.com/alexander-littleton/cadence-api/internal/controllers"
	"github.com/alexander-littleton/cadence-api/internal/controllers/mocks"
	"github.com/alexander-littleton/cadence-api/internal/models"
	"github.com/alexander-littleton/cadence-api/internal/responses"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"io"
	"io/ioutil"
	"net/http"
	"net/http/httptest"
	"net/url"
)

func GetTestGinContext(w *httptest.ResponseRecorder) *gin.Context {
	gin.SetMode(gin.TestMode)

	ctx, _ := gin.CreateTestContext(w)
	ctx.Request = &http.Request{
		Header: make(http.Header),
		URL:    &url.URL{},
	}

	return ctx
}

func MockJsonPost(c *gin.Context, content any) {
	c.Request.Method = "POST"
	c.Request.Header.Set("Content-Type", "application/json")

	jsonbytes, err := json.Marshal(content)
	if err != nil {
		panic(err)
	}

	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonbytes))
}

var _ = Describe("Main", func() {
	var (
		w           *httptest.ResponseRecorder
		ctx         *gin.Context
		ctrl        *gomock.Controller
		userService *mocks.MockUserService
		target      *controllers.UserController
		data        []byte
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		ctx = GetTestGinContext(w)
		ctrl = gomock.NewController(GinkgoT())
		userService = mocks.NewMockUserService(ctrl)
		target = controllers.NewUserController(userService)

	})

	Context("CreateUser", func() {
		var requestBody models.User
		var userResponse responses.UserResponse
		JustBeforeEach(func() {
			MockJsonPost(ctx, requestBody)
			target.CreateUser(ctx)
			res := w.Result()
			defer res.Body.Close()
			data, _ = ioutil.ReadAll(res.Body)
			json.Unmarshal(data, &userResponse)
		})
		Context("the request is valid", func() {
			var newUser models.User
			BeforeEach(func() {
				requestBody = models.User{Email: "test@test.com"}
				newUser = models.User{Id: primitive.NewObjectID(), Email: requestBody.Email}
				userService.EXPECT().CreateUser(ctx, requestBody).Return(newUser, nil)
			})
			It("returns a 201 with a success message and the newly created user", func() {
				output := map[string]interface{}{"data": map[string]interface{}{
					"email": newUser.Email,
					"id":    newUser.Id.Hex(),
				}}
				Expect(userResponse.Data).To(Equal(output))
				Expect(userResponse.Message).To(Equal("success"))
				Expect(w.Code).To(Equal(201))
			})
		})
		Context("new user fails validation", func() {
			BeforeEach(func() {
				requestBody = models.User{Email: ""}
				userService.EXPECT().CreateUser(ctx, requestBody).
					Return(models.User{}, cadence_errors.ValidationErr)
			})
			It("returns a 400 with an error", func() {
				output := map[string]interface{}{"data": cadence_errors.ValidationErr.Error()}
				Expect(userResponse.Data).To(Equal(output))
				Expect(w.Code).To(Equal(400))
			})
		})
		Context("there was an error during processing", func() {
			BeforeEach(func() {
				requestBody = models.User{Email: ""}
				userService.EXPECT().CreateUser(ctx, requestBody).
					Return(models.User{}, errors.New("boom"))
			})
			It("returns a 500 with an error", func() {
				output := map[string]interface{}{"data": "boom"}
				Expect(userResponse.Data).To(Equal(output))
				Expect(w.Code).To(Equal(500))
			})
		})
	})
})
