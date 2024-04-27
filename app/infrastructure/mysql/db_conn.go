package mysql

import (
	"fmt"
	"log"

	"github.com/link-identity/app/config"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DbConn ...
type DbConn struct {
	GormConn *gorm.DB
}

var gormConn *gorm.DB

func init() {
	// TODO : Add once
	dsn := getDSN()
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Println("Failed to connect to database")
		panic(err.Error())
	}
	conn.Exec("CREATE DATABASE IF NOT EXISTS " + config.Values.Database.Name)
	conn.Exec("USE " + config.Values.Database.Name)

	// TODO: enable debug for dev and staging mode
	conn.Logger = logger.Default.LogMode(logger.Info)
	gormConn = conn
	RunMigrations(gormConn)
	log.Println("Successfully connected to database and Ran Migrations")
}

func getDSN() string {
	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/?charset=utf8mb4&parseTime=True&loc=Local",
		config.Values.Database.Username,
		config.Values.Database.Password,
		config.Values.Database.Host,
		config.Values.Database.Port,
	)
	//dsn += config.Values.Database.Name + "?charset=utf8mb4&parseTime=True&loc=Local"
	log.Println(dsn)
	return dsn
}

// NewDBConnection ...
func NewDBConnection() *DbConn {
	return &DbConn{
		GormConn: gormConn,
	}
}
