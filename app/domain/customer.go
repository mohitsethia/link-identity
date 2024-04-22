package domain

import (
	"database/sql"
	"gorm.io/gorm"
	"log"

	"github.com/google/uuid"
)

type Customer struct {
	Model
	CustomerId       string         `json:"user_id,omitempty" gorm:"unique; not null"`
	Email            sql.NullString `db:"email" gorm:"column:email"`
	Phone            sql.NullString `db:"phone" gorm:"column:phone"`
	LinkedID         uint           `json:"linked_id,omitempty"`
	LinkedPrecedence string         `json:"linked_precedence,omitempty" gorm:"not null"`
	Deleted          sql.NullBool   `db:"deleted" gorm:"column:deleted"`
}

func (cus *Customer) TableName() string {
	return "customer"
}
func (cus *Customer) BeforeCreate(scope *gorm.DB) (err error) {
	cus.CustomerId = uuid.New().String()
	log.Println("user generated with id: " + cus.CustomerId)
	return
}
