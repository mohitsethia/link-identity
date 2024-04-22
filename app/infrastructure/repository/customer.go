package repository

import (
	"github.com/link-identity/app/domain"

	"github.com/gin-gonic/gin"
)

type CustomerRepository interface {
	GetCustomerByEmail(ctx *gin.Context, email string) (*domain.Customer, error)
	GetCustomerByPhone(ctx *gin.Context, phone string) (*domain.Customer, error)
	GetAllCustomers(ctx *gin.Context) ([]*domain.Customer, error)
	GetAllSecondaryCustomers(ctx *gin.Context, linkedID uint) ([]*domain.Customer, error)
}
