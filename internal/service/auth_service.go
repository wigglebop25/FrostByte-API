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

func (s *AuthService) Login(username, password string) (*domain.User, string, string, error) {
	user, err := s.repo.FindByUsername(username)
	if err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	var roleNames []string
	for _, role := range user.Roles {
		roleNames = append(roleNames, role.Name)
	}

	// Access Token (15 minutes)
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.UserID,
		"username": user.Username,
		"roles":    roleNames,
		"type":     "access",
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
		"iss":      "frostbyte-api",
	})
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, "", "", err
	}

	// Refresh Token (7 days)
	refreshToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  user.UserID,
		"username": user.Username, // Needed to generate new access token
		"roles":    roleNames,     // Needed to generate new access token
		"type":     "refresh",
		"exp":      time.Now().Add(time.Hour * 24 * 7).Unix(),
		"iss":      "frostbyte-api",
	})
	refreshTokenString, err := refreshToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessTokenString, refreshTokenString, nil
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

	// Verify token type
	if claims["type"] != "refresh" {
		return "", errors.New("invalid token type: expected refresh token")
	}

	userIDFloat, ok := claims["user_id"].(float64)
	if !ok {
		return "", errors.New("user ID not found in token")
	}

	username, _ := claims["username"].(string)
	roles, _ := claims["roles"].([]interface{})

	// Generate new Access Token (15 mins)
	newToken := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id":  uint(userIDFloat),
		"username": username,
		"roles":    roles,
		"type":     "access",
		"exp":      time.Now().Add(time.Minute * 15).Unix(),
		"iss":      "frostbyte-api",
	})

	return newToken.SignedString(s.jwtSecret)
}
