package service

import (
	"frostbyte-api/internal/domain"
	"frostbyte-api/internal/repository"
)

type CategoryService struct {
	repo        *repository.CategoryRepository
	productRepo *repository.ProductRepository
}

func NewCategoryService(repo *repository.CategoryRepository, productRepo *repository.ProductRepository) *CategoryService {
	return &CategoryService{
		repo:        repo,
		productRepo: productRepo,
	}
}

func (s *CategoryService) CreateCategory(category *domain.Category) error {
	return s.repo.Create(category)
}

func (s *CategoryService) GetAllCategories() ([]domain.Category, error) {
	return s.repo.GetAll()
}

func (s *CategoryService) GetCategoryByID(id uint) (*domain.Category, error) {
	return s.repo.FindByID(id)
}

func (s *CategoryService) GetCategoryByName(name string) (*domain.Category, error) {
	return s.repo.FindByName(name)
}

func (s *CategoryService) GetProductByName(name string) (*domain.Product, error) {
	return s.productRepo.FindByName(name)
}

func (s *CategoryService) UpdateCategory(category *domain.Category) error {
	return s.repo.Update(category)
}

func (s *CategoryService) DeleteCategory(id uint) error {
	return s.repo.Delete(id)
}

func (s *CategoryService) AddProductToCategory(categoryID, productID uint) error {
	return s.repo.AddProduct(categoryID, productID)
}

func (s *CategoryService) RemoveProductFromCategory(categoryID, productID uint) error {
	return s.repo.RemoveProduct(categoryID, productID)
}
