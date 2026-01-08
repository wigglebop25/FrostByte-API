package repository

import (
	"frostbyte-api/internal/domain"
	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository(db *gorm.DB) *CategoryRepository {
	return &CategoryRepository{db: db}
}

func (r *CategoryRepository) Create(category *domain.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) GetAll() ([]domain.Category, error) {
	var categories []domain.Category
	err := r.db.Preload("Products").Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) FindByName(name string) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Preload("Products").Where("name = ?", name).First(&category).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) FindByID(id uint) (*domain.Category, error) {
	var category domain.Category
	err := r.db.Preload("Products").First(&category, id).Error
	if err != nil {
		return nil, err
	}
	return &category, nil
}

func (r *CategoryRepository) Update(category *domain.Category) error {
	return r.db.Model(category).Updates(category).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&domain.Category{}, id).Error
}

func (r *CategoryRepository) AddProduct(categoryID, productID uint) error {
	category := domain.Category{CategoryID: categoryID}
	product := domain.Product{ProductID: productID}
	return r.db.Model(&category).Association("Products").Append(&product)
}

func (r *CategoryRepository) RemoveProduct(categoryID, productID uint) error {
	category := domain.Category{CategoryID: categoryID}
	product := domain.Product{ProductID: productID}
	return r.db.Model(&category).Association("Products").Delete(&product)
}
