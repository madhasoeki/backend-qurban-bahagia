package config

import (
	"log"
	"os"
	"qurban/models"
)

func SeedDatabase() {
	var admin models.User
	if DB.Where("username = ?", "admin").First(&admin).Error != nil {
		password := os.Getenv("ADMIN_DEFAULT_PASSWORD")
		if password == "" {
			log.Fatal("ADMIN_DEFAULT_PASSWORD environment variable is not set")
		}

		newAdmin := models.User{
			NamaLengkap: "Administrator Utama",
			Username:    "admin",
			Password:    password,
			Role:        models.RoleAdmin,
		}

		if err := DB.Create(&newAdmin).Error; err != nil {
			log.Fatalf("Failed to seed admin user: %v", err)
		}
		log.Println("Seeder: Admin user created")
	}
}
