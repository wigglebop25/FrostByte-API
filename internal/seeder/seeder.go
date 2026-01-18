package seeder

import (
	"fmt"
	"log"
	"os"
	"time"

	"frostbyte-api/internal/domain"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Seeder handles database seeding with initial data
type Seeder struct {
	db *gorm.DB
}

// NewSeeder creates a new Seeder instance
func NewSeeder(db *gorm.DB) *Seeder {
	return &Seeder{db: db}
}

// SeedIfEmpty checks if the database is empty and seeds it with initial data
func (s *Seeder) SeedIfEmpty() error {
	// Check if users table is empty (excluding roles which are seeded by database.go)
	var userCount int64
	s.db.Model(&domain.User{}).Count(&userCount)

	if userCount > 0 {
		log.Println("Database already has users, skipping seed")
		return nil
	}

	log.Println("Starting database seed...")

	// Seed Roles (should already exist from database.go, but ensure they exist)
	if err := s.seedRoles(); err != nil {
		log.Printf("Error seeding roles: %v", err)
		return err
	}

	// Seed Categories
	if err := s.seedCategories(); err != nil {
		log.Printf("Error seeding categories: %v", err)
		return err
	}

	// Seed Products
	if err := s.seedProducts(); err != nil {
		log.Printf("Error seeding products: %v", err)
		return err
	}

	// Seed Users
	if err := s.seedUsers(); err != nil {
		log.Printf("Error seeding users: %v", err)
		return err
	}

	// Seed Historical Orders
	if err := s.seedHistoricalOrders(); err != nil {
		log.Printf("Error seeding historical orders: %v", err)
		return err
	}

	log.Println("Database seed completed successfully!")
	return nil
}

// seedRoles seeds the roles table
func (s *Seeder) seedRoles() error {
	roles := []domain.Role{
		{Name: "Admin", Permissions: "ADMIN", Description: "Full access"},
		{Name: "Customer", Permissions: "READ,WRITE", Description: "Can view products and place orders"},
		{Name: "Cashier", Permissions: "READ,UPDATE", Description: "Can view orders and update status"},
	}

	for _, r := range roles {
		if err := s.db.FirstOrCreate(&r, domain.Role{Name: r.Name}).Error; err != nil {
			return err
		}
	}

	log.Println("✓ Roles seeded")
	return nil
}

// seedCategories seeds the categories table
func (s *Seeder) seedCategories() error {
	categories := []domain.Category{
		{Name: "Zensai", Description: "Appetizers and Starters"},
		{Name: "Sushi & Sashimi", Description: "Fresh Sushi and Sashimi"},
		{Name: "Menrui", Description: "Noodle Dishes"},
		{Name: "Donburi", Description: "Rice Bowls"},
		{Name: "Kanmi", Description: "Desserts and Sweets"},
	}

	for _, c := range categories {
		if err := s.db.FirstOrCreate(&c, domain.Category{Name: c.Name}).Error; err != nil {
			return err
		}
	}

	log.Println("✓ Categories seeded")
	return nil
}

// seedProducts seeds the products table and associates them with categories
func (s *Seeder) seedProducts() error {
	products := []struct {
		product    domain.Product
		categories []string
	}{
		{
			product: domain.Product{
				Name:            "Pork Gyoza",
				Description:     "Six pieces of pan-seared dumplings served with a soy-vinegar dipping sauce",
				Price:           8.50,
				ProductImageURI: "/images/pork-gyoza.jpg",
			},
			categories: []string{"Zensai"},
		},
		{
			product: domain.Product{
				Name:            "Edamame",
				Description:     "Steamed young soybeans tossed in flakey sea salt",
				Price:           5.00,
				ProductImageURI: "/images/edamame.jpg",
			},
			categories: []string{"Zensai"},
		},
		{
			product: domain.Product{
				Name:            "Shrimp Tempura",
				Description:     "Three pieces of lightly battered, crispy fried shrimp",
				Price:           11.00,
				ProductImageURI: "/images/shrimp-tempura.jpg",
			},
			categories: []string{"Zensai"},
		},
		{
			product: domain.Product{
				Name:            "Maguro Nigiri",
				Description:     "Two pieces of fresh Bluefin tuna over hand-pressed vinegared rice",
				Price:           12.00,
				ProductImageURI: "/images/maguro-nigiri.jpg",
			},
			categories: []string{"Sushi & Sashimi"},
		},
		{
			product: domain.Product{
				Name:            "California Roll",
				Description:     "Classic roll with crab mix, avocado, and cucumber",
				Price:           9.50,
				ProductImageURI: "/images/cali-roll.jpg",
			},
			categories: []string{"Sushi & Sashimi"},
		},
		{
			product: domain.Product{
				Name:            "Salmon Sashimi",
				Description:     "Five thick slices of premium fresh Atlantic salmon",
				Price:           14.50,
				ProductImageURI: "/images/salmon-sashimi.jpg",
			},
			categories: []string{"Sushi & Sashimi"},
		},
		{
			product: domain.Product{
				Name:            "Tonkotsu Ramen",
				Description:     "Rich pork bone broth, chashu pork, bamboo shoots, and a soft-boiled egg",
				Price:           15.99,
				ProductImageURI: "/images/tonkotsu-ramen.jpg",
			},
			categories: []string{"Menrui"},
		},
		{
			product: domain.Product{
				Name:            "Tempura Udon",
				Description:     "Thick wheat noodles in a clear dashi broth topped with shrimp tempura",
				Price:           13.50,
				ProductImageURI: "/images/tempura-udon.jpg",
			},
			categories: []string{"Menrui"},
		},
		{
			product: domain.Product{
				Name:            "Vegetable Yakisoba",
				Description:     "Stir-fried buckwheat noodles with cabbage, carrots, and savory sauce",
				Price:           12.00,
				ProductImageURI: "/images/yakisoba.jpg",
			},
			categories: []string{"Menrui"},
		},
		{
			product: domain.Product{
				Name:            "Gyu-Don",
				Description:     "Thinly sliced beef and onions simmered in a sweet soy dashi over rice",
				Price:           12.50,
				ProductImageURI: "/images/gyudon.jpg",
			},
			categories: []string{"Donburi"},
		},
		{
			product: domain.Product{
				Name:            "Katsu-Don",
				Description:     "Crispy pork cutlet and egg simmered in savory broth over rice",
				Price:           13.50,
				ProductImageURI: "/images/katsudon.jpg",
			},
			categories: []string{"Donburi"},
		},
		{
			product: domain.Product{
				Name:            "Unagi-Don",
				Description:     "Grilled freshwater eel glazed with sweet tare sauce over steamed rice",
				Price:           21.00,
				ProductImageURI: "/images/unagi-don.jpg",
			},
			categories: []string{"Donburi"},
		},
		{
			product: domain.Product{
				Name:            "Matcha Mochi",
				Description:     "Sweet glutinous rice cake filled with premium green tea cream",
				Price:           6.50,
				ProductImageURI: "/images/matcha-mochi.jpg",
			},
			categories: []string{"Kanmi"},
		},
		{
			product: domain.Product{
				Name:            "Taiyaki",
				Description:     "Fish-shaped waffle cake filled with sweet red bean paste",
				Price:           7.00,
				ProductImageURI: "/images/taiyaki.jpg",
			},
			categories: []string{"Kanmi"},
		},
		{
			product: domain.Product{
				Name:            "Black Sesame Ice Cream",
				Description:     "Creamy, nutty, and slightly savory roasted black sesame frozen treat",
				Price:           5.50,
				ProductImageURI: "/images/sesame-ice-cream.jpg",
			},
			categories: []string{"Kanmi"},
		},
	}

	for _, p := range products {
		if err := s.db.FirstOrCreate(&p.product, domain.Product{Name: p.product.Name}).Error; err != nil {
			return err
		}

		// Associate with categories
		for _, catName := range p.categories {
			var category domain.Category
			if err := s.db.Where("name = ?", catName).First(&category).Error; err != nil {
				continue
			}
			s.db.Model(&p.product).Association("Categories").Append(&category)
		}
	}

	log.Println("✓ Products seeded")
	return nil
}

// seedUsers seeds the users table with test users
func (s *Seeder) seedUsers() error {
	// Get passwords from environment variables.
	// CRITICAL: Fail if not set to prevent insecure defaults in production.
	adminPwd := os.Getenv("SEED_ADMIN_PASSWORD")
	if adminPwd == "" {
		return fmt.Errorf("CRITICAL ERROR: SEED_ADMIN_PASSWORD environment variable is not set")
	}

	cashierPwd := os.Getenv("SEED_CASHIER_PASSWORD")
	if cashierPwd == "" {
		return fmt.Errorf("CRITICAL ERROR: SEED_CASHIER_PASSWORD environment variable is not set")
	}

	customerPwd := os.Getenv("SEED_CUSTOMER_PASSWORD")
	if customerPwd == "" {
		return fmt.Errorf("CRITICAL ERROR: SEED_CUSTOMER_PASSWORD environment variable is not set")
	}

	users := []struct {
		username string
		password string
		roleName string
	}{
		{"admin", adminPwd, "Admin"},
		{"cashier", cashierPwd, "Cashier"},
		{"customer", customerPwd, "Customer"},
	}

	for _, u := range users {
		// Hash the password
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(u.password), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		user := domain.User{
			Username:     u.username,
			PasswordHash: string(hashedPassword),
		}

		// Find or create the user
		result := s.db.FirstOrCreate(&user, domain.User{Username: u.username})
		if result.Error != nil {
			return result.Error
		}

		// Assign role to user
		var role domain.Role
		if err := s.db.Where("name = ?", u.roleName).First(&role).Error; err != nil {
			log.Printf("Warning: Role %s not found", u.roleName)
			continue
		}

		// Use raw SQL to insert user role if not already assigned
		s.db.Exec("INSERT IGNORE INTO user_roles (user_id, role_id) VALUES (?, ?)", user.UserID, role.RoleID)
	}

	log.Println("✓ Users seeded")
	return nil
}

// seedHistoricalOrders seeds historical orders for the last 7 days
func (s *Seeder) seedHistoricalOrders() error {
	// Get customer user
	var customerUser domain.User
	if err := s.db.Where("username = ?", "customer").First(&customerUser).Error; err != nil {
		log.Println("Warning: Customer user not found, skipping historical orders")
		return nil
	}

	// Get some products to order
	var products []domain.Product
	if err := s.db.Limit(5).Find(&products).Error; err != nil {
		return err
	}

	if len(products) == 0 {
		log.Println("Warning: No products found, skipping historical orders")
		return nil
	}

	// Generate orders for the last 7 days
	for daysAgo := 6; daysAgo >= 0; daysAgo-- {
		// Skip day 3 for testing gap days
		if daysAgo == 3 {
			continue
		}

		// Create 1-3 orders per day
		numOrders := 1 + (daysAgo % 3)
		for i := 0; i < numOrders; i++ {
			// Create order with some products
			order := domain.Order{
				UserID:      customerUser.UserID,
				Status:      "COMPLETED",
				TotalAmount: float64(15 + (i * 5)),
				CreatedAt:   domain.JSONTime(time.Now().AddDate(0, 0, -daysAgo)),
				UpdatedAt:   domain.JSONTime(time.Now().AddDate(0, 0, -daysAgo)),
			}

			if err := s.db.Create(&order).Error; err != nil {
				return err
			}

			// Add some order products
			product := products[i%len(products)]
			orderProduct := domain.OrderProduct{
				OrderID:   order.OrderID,
				ProductID: product.ProductID,
				Quantity:  1 + (i % 2),
				UnitPrice: product.Price,
				LineTotal: product.Price * float64(1+(i%2)),
			}

			if err := s.db.Create(&orderProduct).Error; err != nil {
				return err
			}
		}
	}

	log.Println("✓ Historical orders seeded")
	return nil
}
