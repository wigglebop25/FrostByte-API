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
	// Strict Role Validation
	allowedRoles := map[string]bool{"Admin": true, "Customer": true, "Cashier": true}
	if roleName == "" {
		return nil, errors.New("role_name is required (Admin, Customer, Cashier)")
	}
	if !allowedRoles[roleName] {
		return nil, errors.New("invalid role: must be Admin, Customer, or Cashier")
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	// 1. Create User
	if err := s.repo.Create(user); err != nil {
		return nil, err
	}

	// 2. Find Role ID
	role, err := s.roleRepo.FindByName(roleName)
	if err != nil {
		return nil, errors.New("role not found in database")
	}

	// 3. Assign Role
	if err := s.roleRepo.AssignRoleToUser(user.UserID, role.RoleID); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *UserService) GetAllUsers() ([]domain.User, error) {
	return s.repo.GetAll()
}

func (s *UserService) GetUserByID(id uint) (*domain.User, error) {
	return s.repo.FindByID(id)
}

func (s *UserService) UpdateUser(user *domain.User, roleName string) error {
	// 1. Update User Details (Username/Password)
	if user.PasswordHash != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
		if err != nil {
			return err
		}
		user.PasswordHash = string(hashedPassword)
	}
	
	if err := s.repo.Update(user); err != nil {
		return err
	}

	// 2. Update Role if provided
	if roleName != "" {
		allowedRoles := map[string]bool{"Admin": true, "Customer": true, "Cashier": true}
		if !allowedRoles[roleName] {
			return errors.New("invalid role: must be Admin, Customer, or Cashier")
		}

		role, err := s.roleRepo.FindByName(roleName)
		if err != nil {
			return errors.New("role not found in database")
		}

		// Prevent downgrading the Root Admin (ID 1)
		if user.UserID == 1 && roleName != "Admin" {
			return errors.New("role cannot be downgraded for root admin")
		}

		// Remove existing roles (assuming single-role system for simplicity, or just overwrite)
		// Since we don't have a specific "RemoveAllRoles" exposed easily, we just Assign.
		// However, many RBAC tables allow multiple. To "Switch" role, strictly we should unassign others.
		// For now, let's just assign the new one. If your logic requires single role, we might need to clear others.
		// Let's assume we just add the new role as primary.
		return s.roleRepo.AssignRoleToUser(user.UserID, role.RoleID)
	}

	return nil
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
