package user

import (
	"errors"

	"golang.org/x/crypto/bcrypt"
)

type Service interface {
	RegisterUser(input RegisterUserInput) (User, error)
	Login(input LoginInput) (User, error)
	IsEmailAvailable(input CheckEmailInput) (bool, error)
	SaveAvatar(id int, fileLocation string) (User, error)
	GetUserByID(id int) (User, error)
}

type service struct {
	repository Repository // repo yg hutruf kecil itu variable punya service sendiri bukan punya Repository
}

func NewService(repository Repository) *service {
	return &service{repository}
}

func (s *service) RegisterUser(input RegisterUserInput) (User, error) {

	user := User{}
	user.Name = input.Name
	user.Email = input.Email
	user.Occupation = input.Occupation
	hash, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.MinCost)
	if err != nil {
		return user, err
	}
	user.PasswordHash = string(hash)
	user.Role = "user"
	newUser, err := s.repository.Save(user)
	if err != nil {
		return newUser, err
	}

	return newUser, nil
}
func (s *service) Login(input LoginInput) (User, error) {
	email := input.Email
	password := input.Password

	user, err := s.repository.FindByEmail(email)
	if err != nil {

		return user, err
	}
	if user.ID == 0 {
		return user, errors.New("User tidak ada")
	}
	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))

	if err != nil {
		return user, err
	}
	return user, nil
}

func (s *service) IsEmailAvailable(input CheckEmailInput) (bool, error) {

	email := input.Email
	user, err := s.repository.FindByEmail(email)

	if err != nil {
		return false, err
	}

	if user.ID == 0 {
		return true, nil
	}
	return false, nil
}

func (s *service) SaveAvatar(id int, fileLocation string) (User, error) {
	// 1. dapatkan user berdasar ID
	// 2. update attribute avatar file name
	// 3. simpan perubahan vatar file name

	user, err := s.repository.FindById(id)
	if err != nil {

		return user, err
	}

	user.AvatarFileName = fileLocation
	updateUser, err := s.repository.UpdateUser(user)
	if err != nil {

		return updateUser, err
	}

	return updateUser, nil
}

func (s *service) GetUserByID(id int) (User, error) {

	user, err := s.repository.FindById(id)
	if err != nil {
		return user, err
	}

	if user.ID == 0 {

		return user, errors.New("user not found with id")

	}
	return user, nil
}
