package mysql

import (
	"github.com/link-identity/app/domain"
	"log"

	"gorm.io/gorm"
)

func RunMigrations(db *gorm.DB) {
	m := []interface{}{
		&domain.Customer{},
	}
	err := db.AutoMigrate(m...)
	if err != nil {
		log.Fatalf("error while connecting to the database %s", err)
	}
	//hashPassword, _ := utilities.HashPassword("admin")
	//db.Create(&models.User{
	//	Username: "admin",
	//	Password: hashPassword,
	//	Role:     "SUPER_ADMIN",
	//})
}
