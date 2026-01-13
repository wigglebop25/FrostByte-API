package service

import (
	"errors"
	"time"

	"frostbyte-api/internal/domain"
	"frostbyte-api/internal/repository"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthService struct {
	repo      *repository.UserRepository
	roleRepo  *repository.RoleRepository
	jwtSecret []byte
}

func NewAuthService(repo *repository.UserRepository, roleRepo *repository.RoleRepository, jwtSecret string) *AuthService {
	return &AuthService{
		repo:      repo,
		roleRepo:  roleRepo,
		jwtSecret: []byte(jwtSecret),
	}
}

func (s *AuthService) Register(username, password string) (*domain.User, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	user := &domain.User{
		Username:     username,
		PasswordHash: string(hashedPassword),
	}

	// Determine Role: If this is the first user, make them Admin
	roleID := uint(2) // Default: Customer
	count, err := s.repo.Count()
	if err == nil && count == 0 {
		roleID = uint(1) // Admin
	}

	if err := s.repo.Create(user); err != nil {
		return nil, err
	}
	
	_ = s.roleRepo.AssignRoleToUser(user.UserID, roleID)

	return user, nil
}

func (s *AuthService) Login(username, password string) (string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return "", errors.New("invalid credentials")
	}

	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.UserID,
		"username": user.Username,
		"sub":      user.Username, // Standard claim for User Identity
		"roles":    roleNames,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iss":      "frostbyte-api",
	})

	tokenString, err := token.SignedString(s.jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.Token, error) {
	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return s.jwtSecret, nil
	})
}

func (s *AuthService) RefreshToken(tokenString string) (string, error) {
	token, err := s.ValidateToken(tokenString)
	if err != nil || !token.Valid {
		return "", errors.New("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", errors.New("invalid token claims")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", errors.New("user ID not found in token")
	}

	// Carry over username, sub, and roles
	username, _ := claims["username"].(string)
	sub, _ := claims["sub"].(string)
	roles, _ := claims["roles"].([]interface{})

	// Generate new token
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  uint(userIDFloat),
		"username": username,
		"sub":      sub,
		"roles":    roles,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
		"iss":      "frostbyte-api",
	})

	return newToken.SignedString(s.jwtSecret)
}
