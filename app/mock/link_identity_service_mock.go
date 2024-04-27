package mock

import (
	"context"

	"github.com/link-identity/app/domain"

	"github.com/stretchr/testify/mock"
)

// LinkIdentityServiceMock ...
type LinkIdentityServiceMock struct {
	mock.Mock
}

// Identify ...
func (m *LinkIdentityServiceMock) Identify(ctx context.Context, email, phone string) ([]*domain.Contact, error) {
	args := m.Called(ctx, email, phone)
	return args.Get(0).([]*domain.Contact), args.Error(1)
}
