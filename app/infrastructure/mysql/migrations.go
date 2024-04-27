package mysql

import (
	"log"

	"github.com/link-identity/app/domain"

	"gorm.io/gorm"
)

// RunMigrations ...
func RunMigrations(db *gorm.DB) {
	m := []interface{}{
		&domain.Contact{},
	}
	err := db.AutoMigrate(m...)
	if err != nil {
		log.Fatalf("error while connecting to the database %s", err)
	}
}
