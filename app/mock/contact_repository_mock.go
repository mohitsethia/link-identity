package mock

import (
	"context"

	"github.com/link-identity/app/domain"

	"github.com/stretchr/testify/mock"
)

// ContactRepositoryMock ...
type ContactRepositoryMock struct {
	mock.Mock
}

// GetContactByEmail ...
func (m *ContactRepositoryMock) GetContactByEmail(
	ctx context.Context,
	email string,
) (*domain.Contact, error) {
	args := m.Called(ctx, email)
	return args.Get(0).(*domain.Contact), args.Error(1)
}

// GetContactByPhone ...
func (m *ContactRepositoryMock) GetContactByPhone(
	ctx context.Context,
	phone string,
) (*domain.Contact, error) {
	args := m.Called(ctx, phone)
	return args.Get(0).(*domain.Contact), args.Error(1)
}

// GetAllContacts ...
func (m *ContactRepositoryMock) GetAllContacts(
	ctx context.Context,
) ([]*domain.Contact, error) {
	args := m.Called(ctx)
	return args.Get(0).([]*domain.Contact), args.Error(1)
}

// GetAllSecondaryContacts ...
func (m *ContactRepositoryMock) GetAllSecondaryContacts(
	ctx context.Context,
	linkedID uint,
) ([]*domain.Contact, error) {
	args := m.Called(ctx, linkedID)
	return args.Get(0).([]*domain.Contact), args.Error(1)
}

// GetPrimaryContactFromLinkedID ...
func (m *ContactRepositoryMock) GetPrimaryContactFromLinkedID(
	ctx context.Context,
	linkedID uint,
) (*domain.Contact, error) {
	args := m.Called(ctx, linkedID)
	return args.Get(0).(*domain.Contact), args.Error(1)
}

// CreateContact ...
func (m *ContactRepositoryMock) CreateContact(
	ctx context.Context,
	contact *domain.Contact,
) (*domain.Contact, error) {
	args := m.Called(ctx, contact)
	return args.Get(0).(*domain.Contact), args.Error(1)
}

// UpdateContact ...
func (m *ContactRepositoryMock) UpdateContact(
	ctx context.Context,
	contact *domain.Contact,
) (*domain.Contact, error) {
	args := m.Called(ctx, contact)
	return args.Get(0).(*domain.Contact), args.Error(1)
}
