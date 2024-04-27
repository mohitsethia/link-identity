package mysql

import (
	"github.com/link-identity/app/domain"
	"log"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	m := []interface{}{
		&domain.Contact{},
	}
	err := db.AutoMigrate(m...)
	if err != nil {
		log.Fatalf("error while connecting to the database %s", err)
	}
}
