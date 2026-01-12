package repository

import (
	"frostbyte-api/internal/domain"
	"gorm.io/gorm"
)

type ProductRepository struct {
	db *gorm.DB
}

func NewProductRepository(db *gorm.DB) *ProductRepository {
	return &ProductRepository{db: db}
}

func (r *ProductRepository) Create(product *domain.Product) error {
	return r.db.Create(product).Error
}

func (r *ProductRepository) GetAll() ([]domain.Product, error) {
	var products []domain.Product
	err := r.db.Preload("Categories").Find(&products).Error
	return products, err
}

func (r *ProductRepository) FindByID(id uint) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Categories").First(&product, id).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) FindByName(name string) (*domain.Product, error) {
	var product domain.Product
	err := r.db.Preload("Categories").Where("name = ?", name).First(&product).Error
	if err != nil {
		return nil, err
	}
	return &product, nil
}

func (r *ProductRepository) Update(product *domain.Product) error {
	return r.db.Transaction(func(tx *gorm.DB) error {
		if err := tx.Model(product).Updates(product).Error; err != nil {
			return err
		}
		
		if product.Categories != nil {
			if err := tx.Model(product).Association("Categories").Replace(product.Categories); err != nil {
				return err
			}
		}
		
		return nil
	})
}

func (r *ProductRepository) Delete(id uint) error {
	var product domain.Product
	if err := r.db.First(&product, id).Error; err != nil {
		return err
	}
	
	// Clear the many-to-many association with categories
	if err := r.db.Model(&product).Association("Categories").Clear(); err != nil {
		return err
	}

	return r.db.Delete(&product).Error
}
