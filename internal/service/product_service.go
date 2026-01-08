package service

import (
	"frostbyte-api/internal/domain"
	"frostbyte-api/internal/repository"
)

type ProductService struct {
	repo         *repository.ProductRepository
	categoryRepo *repository.CategoryRepository
}

func NewProductService(repo *repository.ProductRepository, categoryRepo *repository.CategoryRepository) *ProductService {
	return &ProductService{
		repo:         repo,
		categoryRepo: categoryRepo,
	}
}

func (s *ProductService) CreateProductWithCategories(product *domain.Product, categoryNames []string) error {
	var categories []domain.Category
	for _, name := range categoryNames {
		cat, err := s.categoryRepo.FindByName(name)
		if err == nil {
			categories = append(categories, *cat)
		}
	}
	product.Categories = categories
	return s.repo.Create(product)
}

func (s *ProductService) UpdateProductWithCategories(product *domain.Product, categoryNames []string) error {
	var categories []domain.Category
	for _, name := range categoryNames {
		cat, err := s.categoryRepo.FindByName(name)
		if err == nil {
			categories = append(categories, *cat)
		}
	}
	product.Categories = categories
	return s.repo.Update(product)
}

func (s *ProductService) CreateProduct(product *domain.Product) error {
	return s.repo.Create(product)
}

func (s *ProductService) GetAllProducts() ([]domain.Product, error) {
	return s.repo.GetAll()
}

func (s *ProductService) GetProductByID(id uint) (*domain.Product, error) {
	return s.repo.FindByID(id)
}

func (s *ProductService) UpdateProduct(product *domain.Product) error {
	return s.repo.Update(product)
}

func (s *ProductService) DeleteProduct(id uint) error {
	return s.repo.Delete(id)
}
