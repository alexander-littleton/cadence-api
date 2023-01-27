package api_test

import (
	"bytes"
	"encoding/json"
	"errors"
	"github.com/alexander-littleton/cadence-api/pkg/common/cadence_errors"
	"github.com/alexander-littleton/cadence-api/pkg/user/api"
	"github.com/alexander-littleton/cadence-api/pkg/user/api/mocks"
	"github.com/alexander-littleton/cadence-api/pkg/user/domain"
	"github.com/gin-gonic/gin"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"net/http"
	"net/http/httptest"
)

var _ = Describe("Main", func() {
	var (
		w           *httptest.ResponseRecorder
		router      *gin.Engine
		ctrl        *gomock.Controller
		userService *mocks.MockUserService
		target      api.Controller
	)

	BeforeEach(func() {
		w = httptest.NewRecorder()
		router = gin.New()
		ctrl = gomock.NewController(GinkgoT())
		userService = mocks.NewMockUserService(ctrl)
		target = api.New(userService)
		target.RegisterRoutes(router)

	})

	Context("create new user", func() {
		var newUser domain.User
		JustBeforeEach(func() {
			data, _ := json.Marshal(newUser)
			body := bytes.NewReader(data)
			request, _ := http.NewRequest("POST", "/user", body)
			router.ServeHTTP(w, request)
		})
		Context("the request is valid", func() {
			BeforeEach(func() {
				newUser = domain.User{Email: "test@test.com"}
				createdUser := domain.User{Id: primitive.NewObjectID(), Email: newUser.Email}
				userService.EXPECT().CreateUser(gomock.Any(), mock.MatchedBy(func(u domain.User) bool {
					return u.Email == newUser.Email
				})).Return(createdUser, nil)
			})
			It("returns a 201", func() {
				Expect(w.Code).To(Equal(201))
			})
		})
		Context("new user fails validation", func() {
			BeforeEach(func() {
				newUser = domain.User{Email: ""}
				userService.EXPECT().CreateUser(gomock.Any(), newUser).
					Return(domain.User{}, cadence_errors.ValidationErr)
			})
			It("returns a 400", func() {
				Expect(w.Code).To(Equal(400))
			})
		})
		Context("there was an error during processing", func() {
			BeforeEach(func() {
				newUser = domain.User{Email: ""}
				userService.EXPECT().CreateUser(gomock.Any(), newUser).
					Return(domain.User{}, errors.New("boom"))
			})
			It("returns a 500 with an error", func() {
				Expect(w.Code).To(Equal(500))
			})
		})
	})
})
