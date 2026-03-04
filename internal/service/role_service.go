package service

import (
	"frostbyte-api/internal/domain"
	"frostbyte-api/internal/repository"
)

type RoleService struct {
	repo     *repository.RoleRepository
	userRepo *repository.UserRepository
}

func NewRoleService(repo *repository.RoleRepository, userRepo *repository.UserRepository) *RoleService {
	return &RoleService{repo: repo, userRepo: userRepo}
}

func (s *RoleService) CreateRole(role *domain.Role) error {
	return s.repo.Create(role)
}

func (s *RoleService) GetAllRoles() ([]domain.Role, error) {
	return s.repo.GetAll()
}

func (s *RoleService) GetRoleByID(id uint) (*domain.Role, error) {
	return s.repo.FindByID(id)
}

func (s *RoleService) UpdateRole(role *domain.Role) error {
	return s.repo.Update(role)
}

func (s *RoleService) DeleteRole(id uint) error {
	return s.repo.Delete(id)
}

func (s *RoleService) AddPermission(roleID uint, permission string) error {
	return s.repo.AddPermission(roleID, permission)
}

func (s *RoleService) AssignRoleToUser(userID, roleID uint) error {
	return s.repo.AssignRoleToUser(userID, roleID)
}

func (s *RoleService) AssignRoleToUserByName(username, roleName string) error {
	user, err := s.userRepo.FindByUsername(username)
	if err != nil {
		return err
	}

	role, err := s.repo.FindByName(roleName)
	if err != nil {
		return err
	}

	return s.repo.AssignRoleToUser(user.UserID, role.RoleID)
}
