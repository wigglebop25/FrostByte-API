package service

import (
	"errors"

	"frostbyte-api/internal/domain"
	"frostbyte-api/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo     *repository.UserRepository
	roleRepo *repository.RoleRepository
}

func NewUserService(repo *repository.UserRepository, roleRepo *repository.RoleRepository) *UserService {
	return &UserService{
		repo:     repo,
		roleRepo: roleRepo,
	}
}

func (s *UserService) CreateUser(username, password, roleName string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	
	// Assign Role if provided, otherwise default to Customer (2)
	if roleName != "" {
		// Lookup role by name
		role, err := s.roleRepo.FindByName(roleName)
		if err == nil {
			_ = s.roleRepo.AssignRoleToUser(user.UserID, role.RoleID)
			return user, nil
		}
		// If role not found, fall back to default? Or error?
		// For robustness, let's default to customer if not found or error,
		// or maybe we should return error.
		// Given the prompt "assigning the role during creating", let's try to honor it.
	}
	
	// Default to Customer (2)
	_ = s.roleRepo.AssignRoleToUser(user.UserID, 2)

	return user, nil
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) UpdateUser(user *domain.User) error {
	// Ideally check if user exists first
	return s.repo.Update(user)
}

func (s *UserService) DeleteUser(id uint) error {
	if id == 1 {
		return errors.New("cannot delete root admin user")
	}
	return s.repo.Delete(id)
}

func (s *UserService) SearchUsers(query string) ([]domain.User, error) {
	return s.repo.Search(query)
}
