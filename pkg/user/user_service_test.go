package user_test

import (
	"context"
	"errors"

	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"github.com/stretchr/testify/mock"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/alexander-littleton/cadence-api/pkg/common/cadence_errors"
	"github.com/alexander-littleton/cadence-api/pkg/user"
	"github.com/alexander-littleton/cadence-api/pkg/user/domain"
	mockRepo "github.com/alexander-littleton/cadence-api/pkg/user/repositories/mocks"
)

var _ = Describe("Main", func() {
	var (
		ctrl     *gomock.Controller
		userRepo *mockRepo.MockUserRepository
		target   user.Service
		ctx      context.Context
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		userRepo = mockRepo.NewMockUserRepository(ctrl)
		target = user.New(userRepo)
		ctx = context.TODO()
	})

	Context("createUser", func() {
		var (
			user        domain.User
			createdUser domain.User
			err         error
		)
		BeforeEach(func() {
			user = domain.User{Email: "test@test.com"}
		})
		JustBeforeEach(func() {
			createdUser, err = target.CreateUser(ctx, user)
		})
		Context("the new user is valid", func() {
			BeforeEach(func() {
				userRepo.EXPECT().GetUserByEmail(gomock.Any(), user.Email).Return(domain.User{}, cadence_errors.ErrNotFound)
				userRepo.EXPECT().CreateUser(gomock.Any(), mock.MatchedBy(func(u domain.User) bool {
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
				Expect(createdUser).To(Equal(domain.User{}))
			})
		})
		Context("failed to get existing user with matching email", func() {
			BeforeEach(func() {
				userRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(domain.User{}, errors.New("boom"))
			})
			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("failed to get user by email"))
				Expect(createdUser).To(Equal(domain.User{}))
			})
		})
		Context("a user with matching email already exists", func() {
			BeforeEach(func() {
				userRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(domain.User{Email: user.Email}, nil)
			})
			It("returns a validation Err", func() {
				Expect(err).To(Not(BeNil()))
				Expect(errors.Is(err, cadence_errors.ValidationErr)).To(BeTrue())
				Expect(createdUser).To(Equal(domain.User{}))
			})
		})
		Context("the repository layer returns an error", func() {
			BeforeEach(func() {
				userRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(domain.User{}, cadence_errors.ErrNotFound)
			})
			//It("returns a validation error", func() {
			//	Expect(err).To(Not(BeNil()))
			//	Expect(errors.Is(err, cadence_errors.ValidationErr)).To(BeTrue())
			//	Expect(createdUser).To(Equal(domain.User{}))
			//})
		})
		Context("the repository layer returns an error", func() {
			BeforeEach(func() {
				userRepo.EXPECT().GetUserByEmail(ctx, user.Email).Return(domain.User{}, cadence_errors.ErrNotFound)
				userRepo.EXPECT().CreateUser(ctx, mock.MatchedBy(func(u domain.User) bool {
					return u.Email == user.Email
				})).Return(errors.New("boom"))
			})
			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("failed to create user"))
				Expect(createdUser).To(Equal(domain.User{}))
			})
		})
	})
	Context("GetUserById", func() {
		var (
			userId       primitive.ObjectID
			expectedUser domain.User
			user         domain.User
			err          error
		)
		JustBeforeEach(func() {
			user, err = target.GetUserById(ctx, userId)
		})
		Context("the request is valid", func() {
			BeforeEach(func() {
				userId = primitive.NewObjectID()
				expectedUser = domain.User{Id: userId, Email: "test@test.com"}
				userRepo.EXPECT().GetUserById(ctx, userId).Return(expectedUser, nil)
			})
			It("returns a valid user", func() {
				Expect(err).To(BeNil())
				Expect(expectedUser).To(Equal(user))
			})
		})
		Context("userId is zero", func() {
			BeforeEach(func() {
				expectedUser = domain.User{Email: "test@test.com"}
				userId = expectedUser.Id
			})
			It("returns a validation error ", func() {
				Expect(err).To(Not(BeNil()))
				Expect(errors.Is(err, cadence_errors.ValidationErr)).To(BeTrue())
				Expect(user).To(Equal(domain.User{}))
			})
		})
	})
	Context("GetUserByEmail", func() {
		var (
			user         domain.User
			expectedUser domain.User
			err          error
			email        string
		)
		JustBeforeEach(func() {
			user, err = target.GetUserByEmail(ctx, email)
		})
		Context("the email is valid and a user exists", func() {
			BeforeEach(func() {
				email = "test@test.com"
				expectedUser = domain.User{Id: primitive.NewObjectID(), Email: email}
				userRepo.EXPECT().GetUserByEmail(ctx, email).Return(expectedUser, nil)
			})
			It("returns a user", func() {
				Expect(err).To(BeNil())
				Expect(user).To(Equal(expectedUser))
			})
		})
		Context("the email is invalid", func() {
			BeforeEach(func() {
				email = ""
			})
			It("returns a validation error", func() {
				Expect(err).To(Not(BeNil()))
				Expect(errors.Is(err, cadence_errors.ValidationErr)).To(BeTrue())
				Expect(user).To(Equal(domain.User{}))
			})
		})
		Context("the repository layer returns an error", func() {
			BeforeEach(func() {
				email = "test@test.com"
				expectedUser = domain.User{Id: primitive.NewObjectID(), Email: email}
				userRepo.EXPECT().GetUserByEmail(ctx, email).Return(domain.User{}, errors.New("boom"))
			})
			It("returns an error", func() {
				Expect(err).To(Not(BeNil()))
				Expect(err.Error()).To(ContainSubstring("failed to get user with email"))
				Expect(user).To(Equal(domain.User{}))
			})
		})
	})
})
