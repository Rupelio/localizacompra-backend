package user

import (
	"context"
	"errors"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type userService struct {
	repo Repository
}

func NewService(r Repository) Service {
	return &userService{
		repo: r,
	}
}

func (s *userService) Create(ctx context.Context, user User) (User, error) {
	if user.Name == "" {
		return User{}, errors.New("o nome de usuário não pode ser vazio")
	}
	if len(user.PasswordHash) < 8 {
		return User{}, errors.New("a senha precisa ter mais de 8 digitos")
	}

	_, err := s.repo.GetByEmail(ctx, user.Email)
	if err == nil {
		return User{}, errors.New("o email já está em uso")
	}

	if !errors.Is(err, ErrUserNotFound) {
		return User{}, err
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		return User{}, err
	}

	user.PasswordHash = string(hashedPassword)

	return s.repo.Create(ctx, user)
}

var jwtSecret = []byte("sua-chave-super-secreta")

func (s *userService) Login(ctx context.Context, email, password string) (string, error) {
	user, err := s.repo.GetByEmail(ctx, email)
	if err != nil {
		return "", errors.New("credenciais inválidas")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return "", errors.New("credenciais inválidas")
	}

	claims := jwt.MapClaims{
		"sub":  user.ID,
		"role": user.Role,
		"exp":  time.Now().Add(time.Hour * 24).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(jwtSecret)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (s *userService) GetByID(ctx context.Context, id int64) (User, error) {
	return s.repo.GetByID(ctx, id)
}
