package usecase

import (
	"github.com/hardyantz/data-encryption/pkg/member-aes/domain"
	"github.com/hardyantz/data-encryption/pkg/member-aes/repository"
)

type MemberUseCase interface {
	Create(member *domain.Member) error
	GetAll(params domain.Parameters) ([]domain.Member, error)
	GetOne(id string) (*domain.Member, error)
	GetOneEmail(email string) (*domain.Member, error)
	Update(id string, member *domain.Member) error
	Delete(id string) error
	GetEmailLogin(email string) (*domain.Member, error)
}

type memberUseCase struct {
	repository repository.MemberRepository
}

func NewMemberUseCase(repo repository.MemberRepository) MemberUseCase {
	return &memberUseCase{repo}
}

func (u *memberUseCase) Create(member *domain.Member) error {
	if err := u.repository.Create(member); err != nil {
		return err
	}
	return nil
}

func (u *memberUseCase) GetAll(params domain.Parameters) ([]domain.Member, error) {
	members, err := u.repository.GetAll(params)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (u *memberUseCase) GetOne(id string) (*domain.Member, error) {
	return u.repository.GetOne(id)
}

func (u *memberUseCase) GetOneEmail(email string) (*domain.Member, error) {
	return u.repository.GetOneEmail(email)
}

func (u *memberUseCase) Update(id string, member *domain.Member) error {
	return u.repository.Update(id, member)
}

func (u *memberUseCase) Delete(id string) error {
	return u.repository.Delete(id)
}

func (u *memberUseCase) GetEmailLogin(email string) (*domain.Member, error) {
	return u.repository.GetEmailLogin(email)
}
