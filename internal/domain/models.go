package domain

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	UserID       uint           `gorm:"primaryKey;column:user_id" json:"user_id"`
	Username     string         `gorm:"size:50;not null;unique" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Roles        []Role         `gorm:"many2many:user_roles;joinForeignKey:user_id;joinReferences:role_id" json:"roles,omitempty"`
	Orders       []Order        `json:"orders,omitempty"`
	CreatedAt    time.Time      `json:"created_at"`
	UpdatedAt    time.Time      `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Role struct {
	RoleID      uint      `gorm:"primaryKey;column:role_id" json:"role_id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name"`
	Permissions string    `gorm:"type:text" json:"permissions"` // SET type in MySQL, using string for simplicity
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Category struct {
	CategoryID  uint      `gorm:"primaryKey;column:category_id" json:"category_id"`
	Name        string    `gorm:"size:255;not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Products    []Product `gorm:"many2many:product_categories;joinForeignKey:category_id;joinReferences:product_id" json:"products,omitempty"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

type Product struct {
	ProductID       uint       `gorm:"primaryKey;column:product_id" json:"product_id"`
	Name            string     `gorm:"size:100;not null;unique" json:"name"`
	ProductImageURI string     `gorm:"size:255;column:product_image_uri" json:"product_image_uri"`
	Description     string     `gorm:"type:text" json:"description"`
	Price           float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Categories      []Category `gorm:"many2many:product_categories;joinForeignKey:product_id;joinReferences:category_id" json:"categories,omitempty"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

type Order struct {
	OrderID     uint           `gorm:"primaryKey;column:order_id" json:"order_id"`
	UserID      uint           `gorm:"column:user_id;not null" json:"user_id"`
	User        User           `json:"user,omitempty"`
	TotalAmount float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      string         `gorm:"size:50;default:'pending'" json:"status"`
	Products    []OrderProduct `gorm:"foreignKey:OrderID" json:"products,omitempty"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

type OrderProduct struct {
	OrderID   uint    `gorm:"primaryKey;column:order_id" json:"order_id"`
	ProductID uint    `gorm:"primaryKey;column:product_id" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int     `gorm:"not null;default:1" json:"quantity"`
	UnitPrice float64 `gorm:"column:unit_price;type:decimal(10,2);not null" json:"unit_price"`
	LineTotal float64 `gorm:"column:line_total;type:decimal(10,2);->;<-:false" json:"line_total"` // Read-only
}