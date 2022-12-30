package user_service_test

import (
	"context"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"internal/common/cadence_errors"
	"internal/models"
	"internal/services"
	"internal/services/mocks"
)

var _ = Describe("Main", func() {
	var (
		ctrl     *gomock.Controller
		userRepo *mocks.MockUserRepository
		target   *user_service.UserService
		ctx      context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		userRepo = mocks.NewMockUserRepository(ctrl)
		target = user_service.NewUserService(userRepo)
		ctx = context.TODO()
	})

	Context("CreateUser", func() {
		var user models.User
		var createdUser models.User
		var err error
		JustBeforeEach(func() {
			createdUser, err = target.CreateUser(ctx, user)
		})
		Context("the new user is valid", func() {
			BeforeEach(func() {
				user = models.User{Email: "test@test.com"}
				userRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(models.User{}, cadence_errors.ErrNotFound)
				userRepo.EXPECT().CreateUser(ctx, user).Return(nil)
			})
			It("returns the created user", func() {
				Expect(err).To(BeNil())
				Expect(createdUser.Email).To(Equal(user.Email))
			})
		})
	})
})
