package service

import (
	"bilibili-clone/internal/user/model"
	"bilibili-clone/internal/user/repository"
	"bilibili-clone/pkg/utils"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

type UserService struct {
	repo *repository.UserRepository
}

func NewUserService() *UserService {
	return &UserService{
		repo: repository.NewUserRepository(),
	}
}

func (s *UserService) Register(email, password, nickname string) error {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return errors.New("email is required")
	}

	// 如果没有提供nickname，使用email前缀作为默认昵称
	if strings.TrimSpace(nickname) == "" {
		emailParts := strings.Split(email, "@")
		if len(emailParts) > 0 {
			nickname = "User_" + emailParts[0]
		} else {
			nickname = "User"
		}
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	user := &model.User{
		Email:    email,
		Password: string(hashedPassword),
		Nickname: nickname,
		Avatar:   "",
	}

	return s.repo.Create(user)
}

func (s *UserService) Login(email, password string) (string, *model.User, error) {
	email = strings.ToLower(strings.TrimSpace(email))
	if email == "" {
		return "", nil, errors.New("invalid credentials")
	}

	user, err := s.repo.FindByEmail(email)
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", nil, errors.New("invalid credentials")
	}

	token, tokenErr := utils.GenerateToken(user.ID, user.Nickname, user.Avatar)
	if tokenErr != nil {
		return "", nil, tokenErr
	}

	return token, user, nil
}

func (s *UserService) UpdateUser(userID uint, nickname, avatar string) error {
	user, err := s.repo.FindByID(userID)
	if err != nil {
		return err
	}

	if nickname != "" {
		user.Nickname = nickname
	}
	if avatar != "" {
		user.Avatar = avatar
	}

	return s.repo.Update(user)
}

func (s *UserService) GetUsersByIDs(ids []uint) ([]model.User, error) {
	return s.repo.FindByIDs(ids)
}

func (s *UserService) GetUserByID(userID uint) (*model.User, error) {
        return s.repo.FindByID(userID)
}
