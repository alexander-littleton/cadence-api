package userservice_test

import (
	"context"
	"errors"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"internal/common/cadence_errors"
	"internal/models"
	"internal/services"
	"internal/services/mocks"
)

var _ = Describe("Main", func() {
	var (
		ctrl     *gomock.Controller
		userRepo *mocks.MockUserRepository
		target   *userservice.UserService
		ctx      context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		userRepo = mocks.NewMockUserRepository(ctrl)
		target = userservice.NewUserService(userRepo)
		ctx = context.TODO()
	})

	Context("CreateUser", func() {
		var (
			user        models.User
			createdUser models.User
			err         error
		)
		BeforeEach(func() {
			user = models.User{Email: "test@test.com"}
		})
		JustBeforeEach(func() {
			createdUser, err = target.CreateUser(ctx, user)
		})
		Context("the new user is valid", func() {
			BeforeEach(func() {
				userRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(models.User{}, cadence_errors.ErrNotFound)
				userRepo.EXPECT().CreateUser(ctx, mock.MatchedBy(func(u models.User) bool {
					return u.Email == user.Email
				})).Return(nil)
			})
			It("returns the created user", func() {
				Expect(err).To(BeNil())
				Expect(createdUser.Email).To(Equal(user.Email))
			})
		})
		Context("the user already has an object id", func() {
			BeforeEach(func() {
				user.Id = primitive.NewObjectID()
			})
			It("returns a validation Err", func() {
				Expect(err).To(Not(BeNil()))
				Expect(errors.Is(err, cadence_errors.ValidationErr)).To(BeTrue())
			})
		})
		Context("failed to get existing user with matching email", func() {
			BeforeEach(func() {
				userRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(models.User{}, errors.New("boom"))
			})
			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("failed to get user by email"))
			})
		})
		Context("a user with matching email already exists", func() {
			BeforeEach(func() {
				userRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(models.User{Email: user.Email}, nil)
			})
			It("returns a validation Err", func() {
				Expect(err).To(Not(BeNil()))
				Expect(errors.Is(err, cadence_errors.ValidationErr)).To(BeTrue())
			})
		})
		//Cannot pass until we have more fields
		//Context("a user is missing required fields", func() {
		//	BeforeEach(func() {
		//		user = models.User{}
		//		userRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(models.User{Email: user.Email}, nil)
		//	})
		//	It("returns a validation Err", func() {
		//		Expect(err).To(Not(BeNil()))
		//		Expect(errors.Is(err, cadence_errors.ValidationErr)).To(BeTrue())
		//	})
		//})
		Context("the repository layer returns an error", func() {
			BeforeEach(func() {
				userRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(models.User{}, cadence_errors.ErrNotFound)
				userRepo.EXPECT().CreateUser(ctx, mock.MatchedBy(func(u models.User) bool {
					return u.Email == user.Email
				})).Return(errors.New("boom"))
			})
			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("failed to create user"))
			})
		})
	})
})
