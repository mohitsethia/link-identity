package domain

import (
	"database/sql"
)

// Contact ...
type Contact struct {
	Model
	ContactID        uint           `json:"contact_id,omitempty" gorm:"primaryKey; unique; not null; autoIncrement"`
	Email            sql.NullString `db:"email" gorm:"column:email"`
	Phone            sql.NullString `db:"phone" gorm:"column:phone"`
	LinkedID         uint           `json:"linked_id,omitempty"`
	LinkedPrecedence string         `json:"linked_precedence,omitempty" gorm:"not null" default:"primary"`
	Deleted          sql.NullBool   `db:"deleted" gorm:"column:deleted"`
}

// TableName ...
func (c *Contact) TableName() string {
	return "contact"
}
