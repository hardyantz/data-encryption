package usecase

import (
	"github.com/hardyantz/data-encryption/pkg/member/domain"
	repository "github.com/hardyantz/data-encryption/pkg/member/repository"
)

type MemberImplementation interface {
	Create(member *domain.Member) error
	GetAll(params domain.Parameters) ([]domain.Member, error)
	GetOne(id string) (domain.Member, error)
	GetOneEmail(email string) (domain.Member, error)
	Update(id string, member *domain.Member) error
	Delete(id string) error
	GetEmailLogin(email string) (domain.Member, error)
}

type MemberUseCase struct {
	repository *repository.MemberRepositoryMongo
}

func NewMemberImplementation(repo *repository.MemberRepositoryMongo) *MemberUseCase {
	return &MemberUseCase{repo}
}

func (u *MemberUseCase) Create(member *domain.Member) error {
	if err := u.repository.Create(member); err != nil {
		return err
	}
	return nil
}

func (u *MemberUseCase) GetAll(params domain.Parameters) ([]domain.Member, error) {
	members, err := u.repository.GetAll(params)
	if err != nil {
		return nil, err
	}

	return members, nil
}

func (u *MemberUseCase) GetOne(id string) (*domain.Member, error) {
	return u.repository.GetOne(id)
}

func (u *MemberUseCase) GetOneEmail(email string) (*domain.Member, error) {
	return u.repository.GetOneEmail(email)
}

func (u *MemberUseCase) Update(id string, member *domain.Member) error {
	return u.repository.Update(id, member)
}

func (u *MemberUseCase) Delete(id string) error {
	return u.repository.Delete(id)
}

func (u *MemberUseCase) GetEmailLogin(email string) (*domain.Member, error) {
	return u.repository.GetEmailLogin(email)
}
