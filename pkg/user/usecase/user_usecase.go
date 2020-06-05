package usecase

import (
	"github.com/hardyantz/data-encryption/pkg/user/domain"
	"github.com/hardyantz/data-encryption/pkg/user/repository"
)

type UserImplementation interface {
	Create(user *domain.User) error
	GetAll(params domain.Parameters) ([]domain.User, error)
	GetOne(id string) (domain.User, error)
	GetOneEmail(email string) (domain.User, error)
	Update(id string, user *domain.User) error
	Delete(id string) error
}

type UserUseCase struct {
	repository *repository.UserRepository
}

func NewUserImplementation(repo *repository.UserRepository) *UserUseCase {
	return &UserUseCase{repo}
}

func (u *UserUseCase) Create(user *domain.User) error {
	if err := u.repository.Create(user); err != nil {
		return err
	}
	return nil
}

func (u *UserUseCase) GetAll(params domain.Parameters) ([]domain.User, error) {
	users, err := u.repository.GetAll(params)
	if err != nil {
		return nil, err
	}

	return users, nil
}

func (u *UserUseCase) GetOne(id string) (domain.User, error) {
	return u.repository.GetOne(id)
}

func (u *UserUseCase) GetOneEmail(email string) (domain.User, error) {
	return u.repository.GetOneEmail(email)
}

func (u *UserUseCase) Update(id string, user *domain.User) error {
	return u.repository.Update(id, user)
}

func (u *UserUseCase) Delete(id string) error {
	return u.repository.Delete(id)
}