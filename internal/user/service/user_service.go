package service

import (
"bilibili-clone/internal/user/model"
"bilibili-clone/internal/user/repository"
"bilibili-clone/pkg/utils"
"errors"

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

func (s *UserService) Register(username, password string) error {
hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
if err != nil {
return err
}

user := &model.User{
Username: username,
Password: string(hashedPassword),
Nickname: "User" + username, // Default nickname
Avatar:   "",                // Default avatar
}

return s.repo.Create(user)
}

func (s *UserService) Login(username, password string) (string, error) {
user, err := s.repo.FindByUsername(username)
if err != nil {
return "", errors.New("invalid credentials")
}

err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
if err != nil {
return "", errors.New("invalid credentials")
}

return utils.GenerateToken(user.ID, user.Nickname, user.Avatar)
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
