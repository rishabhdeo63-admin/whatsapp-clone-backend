package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/rishabhdeo60-dev/whatsapp-clone/internal/dao"
	"github.com/rishabhdeo60-dev/whatsapp-clone/internal/model"
	"github.com/rishabhdeo60-dev/whatsapp-clone/internal/service"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockUserRepository struct{ mock.Mock }

func (m *MockUserRepository) Create(ctx context.Context, user *model.User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) FindByMobileNumber(ctx context.Context, mobileNumber uint64) (*model.User, error) {
	args := m.Called(mobileNumber)
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserRepository) FindByUsername(ctx context.Context, username string) (*model.User, error) {
	args := m.Called(username)
	return args.Get(0).(*model.User), args.Error(1)
}
func (m *MockUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	args := m.Called(email)
	return args.Get(0).(*model.User), args.Error(1)
}

func (m *MockUserRepository) FindByID(ctx context.Context, userID int64) (*model.User, error) {
	args := m.Called(userID)
	return args.Get(0).(*model.User), args.Error(1)
}

type MockContactRepository struct{ mock.Mock }

func (m *MockContactRepository) AddContact(ctx context.Context, userID, contactID int64) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

func (m *MockContactRepository) GetContacts(ctx context.Context, userID int64) ([]*dao.ContactDAO, error) {
	args := m.Called(ctx, userID)
	return args.Get(0).([]*dao.ContactDAO), args.Error(1)
}

func (m *MockContactRepository) RemoveContact(ctx context.Context, userID, contactID int64) error {
	args := m.Called(ctx, userID, contactID)
	return args.Error(0)
}

// ------------------Tests----------------

func TestAddContact_Success(t *testing.T) {
	uRepo := new(MockUserRepository)
	cRepo := new(MockContactRepository)
	service := service.NewContactService(cRepo, uRepo)

	uRepo.On("FindByID", 2).Return(&model.User{ID: 2}, nil)
	cRepo.On("AddContact", context.Background(), 1, 2).Return(nil)

	err := service.AddContact(1, 2)
	assert.NoError(t, err)
}

func TestAddContact_SelfAddError(t *testing.T) {
	uRepo := new(MockUserRepository)
	cRepo := new(MockContactRepository)
	service := service.NewContactService(cRepo, uRepo)

	err := service.AddContact(1, 1)
	assert.Error(t, err)
	assert.Equal(t, "cannot add yourself as a contact", err.Error())
}

func TestAddContact_UserNotFound(t *testing.T) {
	uRepo := new(MockUserRepository)
	cRepo := new(MockContactRepository)
	service := service.NewContactService(cRepo, uRepo)

	uRepo.On("FindByID", 2).Return(&model.User{}, errors.New("not found"))

	err := service.AddContact(1, 2)
	assert.Error(t, err)
}
