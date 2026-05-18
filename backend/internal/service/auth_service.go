package service

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stockapp/backend/internal/model"
	"github.com/stockapp/backend/internal/repository"
	"github.com/stockapp/backend/internal/utils"
)

type AuthService struct {
	jwtSecret string
}

func NewAuthService(secret string) *AuthService {
	return &AuthService{jwtSecret: secret}
}

func (s *AuthService) Register(req *model.RegisterRequest) (*model.AuthResponse, error) {
	// Check if email exists
	exists, err := repository.EmailExists(req.Email)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errors.New("el email ya está registrado")
	}

	// Hash password
	hash, err := utils.HashPassword(req.Password)
	if err != nil {
		return nil, err
	}

	// Create user with default "usuario" role
	usuario, err := repository.CreateUsuario(req.Email, hash, req.Nombre, model.RoleUsuario)
	if err != nil {
		return nil, err
	}

	// Generate token
	token, err := s.generateToken(usuario)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User:  usuario,
		Token: token,
	}, nil
}

func (s *AuthService) Login(req *model.LoginRequest) (*model.AuthResponse, error) {
	// Get user by email
	usuario, err := repository.GetUsuarioByEmail(req.Email)
	if err != nil {
		return nil, errors.New("email o contraseña incorrectos")
	}

	// Check password
	if !utils.CheckPassword(req.Password, usuario.PasswordHash) {
		return nil, errors.New("email o contraseña incorrectos")
	}

	// Generate token
	token, err := s.generateToken(usuario)
	if err != nil {
		return nil, err
	}

	return &model.AuthResponse{
		User:  usuario,
		Token: token,
	}, nil
}

func (s *AuthService) generateToken(usuario *model.Usuario) (string, error) {
	claims := jwt.MapClaims{
		"sub": usuario.ID,
		"email": usuario.Email,
		"role": usuario.Rol,
		"exp": time.Now().Add(24 * time.Hour * 7).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *AuthService) ValidateToken(tokenString string) (*jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return &claims, nil
	}

	return nil, errors.New("token inválido")
}

func (s *AuthService) GetUserByID(id string) (*model.Usuario, error) {
	return repository.GetUsuarioByID(id)
}