package database

import (
	"log"

	"frostbyte-api/internal/domain"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func Connect(dsn string) {
	var err error
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}

	log.Println("Database connected successfully!")

	// Disable foreign key checks for migration
	DB.Exec("SET FOREIGN_KEY_CHECKS=0")

	// Auto-Migrate models
	err = DB.AutoMigrate(
		&domain.Role{},
		&domain.User{},
		&domain.Category{},
		&domain.Product{},
		&domain.Order{},
		&domain.OrderProduct{},
	)
	if err != nil {
		log.Fatalf("Failed to migrate database: %v", err)
	}

	// Seed Roles
	roles := []domain.Role{
		{Name: "Admin", Permissions: "READ,WRITE,DELETE,UPDATE", Description: "Full access"},
		{Name: "Customer", Permissions: "READ,WRITE", Description: "Can view products and place orders"},
		{Name: "Cashier", Permissions: "READ,UPDATE", Description: "Can view orders and update status"},
	}

	for _, r := range roles {
		DB.FirstOrCreate(&r, domain.Role{Name: r.Name})
	}

	// Re-enable foreign key checks
	DB.Exec("SET FOREIGN_KEY_CHECKS=1")

	log.Println("Database migration completed!")
}
