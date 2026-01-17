package domain

import (
	"database/sql/driver"
	"fmt"
	"time"

	"gorm.io/gorm"
)

// JSONTime is a custom type for human-readable time formatting in JSON
type JSONTime time.Time

func (t JSONTime) MarshalJSON() ([]byte, error) {
	loc, err := time.LoadLocation("Asia/Manila")
	if err != nil {
		fmt.Printf("Error loading location 'Asia/Manila': %v\n", err) // Log the error to stdout
		loc = time.UTC
	}
	formatted := time.Time(t).In(loc).Format("January 02, 2006 03:04 PM")
	stamp := fmt.Sprintf("\"%s\"", formatted)
	return []byte(stamp), nil
}

// Value implements the driver.Valuer interface
func (t JSONTime) Value() (driver.Value, error) {
	return time.Time(t), nil
}

// Scan implements the sql.Scanner interface
func (t *JSONTime) Scan(value interface{}) error {
	if value == nil {
		*t = JSONTime(time.Time{})
		return nil
	}
	if v, ok := value.(time.Time); ok {
		*t = JSONTime(v)
		return nil
	}
	return fmt.Errorf("cannot scan type %T into JSONTime", value)
}

type User struct {
	UserID       uint           `gorm:"primaryKey;column:user_id" json:"user_id"`
	Username     string         `gorm:"size:50;not null;unique" json:"username"`
	PasswordHash string         `gorm:"size:255;not null" json:"-"`
	Roles        []Role         `gorm:"many2many:user_roles;joinForeignKey:user_id;joinReferences:role_id" json:"roles,omitempty"`
	Orders       []Order        `json:"orders,omitempty"`
	CreatedAt    JSONTime       `json:"created_at"`
	UpdatedAt    JSONTime       `json:"updated_at"`
	DeletedAt    gorm.DeletedAt `gorm:"index" json:"-"`
}

type Role struct {
	RoleID      uint      `gorm:"primaryKey;column:role_id" json:"role_id"`
	Name        string    `gorm:"size:50;not null;unique" json:"name"`
	Permissions string    `gorm:"type:text" json:"permissions"`
	Description string    `gorm:"type:text" json:"description"`
	CreatedAt   JSONTime  `json:"created_at"`
	UpdatedAt   JSONTime  `json:"updated_at"`
}

type Category struct {
	CategoryID  uint      `gorm:"primaryKey;column:category_id" json:"category_id"`
	Name        string    `gorm:"size:255;not null;unique" json:"name"`
	Description string    `gorm:"type:text" json:"description"`
	Products    []Product `gorm:"many2many:product_categories;joinForeignKey:category_id;joinReferences:product_id" json:"products,omitempty"`
	CreatedAt   JSONTime  `json:"created_at"`
	UpdatedAt   JSONTime  `json:"updated_at"`
}

type Product struct {
	ProductID       uint       `gorm:"primaryKey;column:product_id" json:"product_id"`
	Name            string     `gorm:"size:100;not null;unique" json:"name"`
	ProductImageURI string     `gorm:"size:255;column:product_image_uri" json:"product_image_uri"`
	Description     string     `gorm:"type:text" json:"description"`
	Price           float64    `gorm:"type:decimal(10,2);not null" json:"price"`
	Categories      []Category `gorm:"many2many:product_categories;joinForeignKey:product_id;joinReferences:category_id" json:"categories,omitempty"`
	CreatedAt       JSONTime   `json:"created_at"`
	UpdatedAt       JSONTime   `json:"updated_at"`
}

type Order struct {
	OrderID     uint           `gorm:"primaryKey;column:order_id" json:"order_id"`
	UserID      uint           `gorm:"column:user_id;not null" json:"user_id"`
	User        User           `json:"user,omitempty"`
	TotalAmount float64        `gorm:"type:decimal(10,2);not null" json:"total_amount"`
	Status      string         `gorm:"size:50;default:'PENDING'" json:"status"`
	Products    []OrderProduct `gorm:"foreignKey:OrderID" json:"products,omitempty"`
	CreatedAt   JSONTime       `json:"created_at"`
	UpdatedAt   JSONTime       `json:"updated_at"`
}

type OrderProduct struct {
	OrderID   uint    `gorm:"primaryKey;column:order_id" json:"order_id"`
	ProductID uint    `gorm:"primaryKey;column:product_id" json:"product_id"`
	Product   Product `gorm:"foreignKey:ProductID" json:"product,omitempty"`
	Quantity  int     `gorm:"not null;default:1" json:"quantity"`
	UnitPrice float64 `gorm:"column:unit_price;type:decimal(10,2);not null" json:"unit_price"`
	LineTotal float64 `gorm:"column:line_total;type:decimal(10,2);->;<-:false" json:"line_total"`
}
