package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hardyantz/data-encryption/config"
	"github.com/hardyantz/data-encryption/helpers"
	"github.com/hardyantz/data-encryption/pkg/member-aes/domain"
)

type MemberRepository interface {
	Create(member *domain.Member) error
	GetAll(domain.Parameters) ([]domain.Member, error)
	GetOne(id string) (*domain.Member, error)
	GetOneEmail(email string) (*domain.Member, error)
	Update(id string, member *domain.Member) error
	Delete(id string) error
	GetEmailLogin(email string) (*domain.Member, error)
}

type memberRepository struct {
	db    *sql.DB
	conf  *config.Config
	redis *config.Redis
}

func NewMemberRepository(db *sql.DB, conf *config.Config, redis *config.Redis) MemberRepository {
	return &memberRepository{db, conf, redis}
}

func (r *memberRepository) Create(member *domain.Member) error {

	member.GenID()
	encryptedEmail, err := helpers.Encrypt(member.Email, config.PassPhrase)
	if err != nil {
		return err
	}
	query := fmt.Sprintf(`INSERT INTO members_aes(
	user_id, first_name, last_name, email, phone, password, email_encrypted) 
	VALUES ('%s', '%s', '%s', '%s', '%s', '%s', '%s')`,
		member.ID, member.FirstName, member.LastName, member.Email, member.Phone, member.Password, encryptedEmail)

	_, err = r.db.Exec(query)
	if err != nil {
		return err
	}

	member.EmailEncrypted = encryptedEmail

	return nil
}

func (r *memberRepository) GetAll(params domain.Parameters) ([]domain.Member, error) {
	var (
		members []domain.Member
		err     error
		sort    string
	)

	q := fmt.Sprintf(`select "user_id", "first_name", "last_name", "email", "phone", "password", CAST(COALESCE(NULLIF("email_encrypted", ''), '') as VARCHAR) as email_encrypted from members_aes`)

	sort = "asc"

	if len(params.OrderBy) > 0 && len(params.Sort) > 0 {
		if helpers.StringInSlice(strings.ToLower(params.Sort), []string{"asc", "desc"}) {
			sort = params.Sort
		}
		q += fmt.Sprintf(` ORDER BY %s %s`, params.OrderBy, sort)
	}

	if len(params.Limit) > 0 {
		q += fmt.Sprintf(" LIMIT %s", params.Limit)
	}

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var member domain.Member
		if err := rows.Scan(&member.ID, &member.FirstName, &member.LastName, &member.Email, &member.Phone, &member.Password, &member.EmailEncrypted); err != nil {
			return nil, err
		}

		member.EmailEncrypted, _ = helpers.Decrypt(member.EmailEncrypted, config.PassPhrase)

		members = append(members, member)
	}

	return members, err
}

func (r *memberRepository) GetOne(id string) (*domain.Member, error) {
	member := new(domain.Member)

	q := fmt.Sprintf(`select "user_id", "first_name", "last_name", "email", "phone", "password", CAST(COALESCE(NULLIF("email_encrypted", ''), '') as VARCHAR) from members_aes where "user_id" = '%s'`, id)

	err := r.db.QueryRow(q).Scan(&member.ID, &member.FirstName, &member.LastName, &member.Email, &member.Phone, &member.Password, &member.EmailEncrypted)

	if err != nil {
		return member, err
	}

	member.EmailEncrypted, _ = helpers.Decrypt(member.EmailEncrypted, config.PassPhrase)

	return member, err
}

func (r *memberRepository) GetOneEmail(email string) (*domain.Member, error) {
	member := new(domain.Member)

	emailEncrypted, _ := helpers.Encrypt(email, config.PassPhrase)
	q := fmt.Sprintf(`select "user_id", "first_name", "last_name", "email", "phone", "password",CAST(COALESCE(NULLIF("email_encrypted", ''), '') as VARCHAR) from members_aes where email_encrypted = '%s'`, emailEncrypted)

	err := r.db.QueryRow(q).Scan(&member.ID, &member.FirstName, &member.LastName, &member.Email, &member.Phone, &member.Password, &member.EmailEncrypted)

	member.EmailEncrypted, _ = helpers.Decrypt(member.EmailEncrypted, config.PassPhrase)

	return member, err
}

func (r *memberRepository) Update(id string, member *domain.Member) error {

	getMember, _ := r.GetOne(id)

	if member.Email != "" {
		emailEncrypted, _ := helpers.Decrypt(member.EmailEncrypted, config.PassPhrase)
		getMember.EmailEncrypted = emailEncrypted
		getMember.EmailEncrypted, _ = helpers.Encrypt(member.Email, config.PassPhrase)
		getMember.Email = member.Email
	}

	query := fmt.Sprintf(`UPDATE members_aes SET
	"first_name" = '%s', 
	"last_name" = '%s', 
	"phone" = '%s', 
	"password" = '%s',
	"email" = '%s',
	"email_encrypted" = '%s'
	WHERE "user_id" = '%s'`,
		member.FirstName, member.LastName, member.Phone, member.Password, getMember.Email, getMember.EmailEncrypted, id)

	res, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	if _, err := res.RowsAffected(); err != nil {
		return err
	}

	member.EmailEncrypted = getMember.EmailEncrypted
	member.Email = getMember.Email

	return nil
}

func (r *memberRepository) Delete(id string) error {

	r.db.QueryRow(`DELETE FROM members_aes WHERE user_id = '%s'`, id)

	return nil
}

func (r *memberRepository) GetEmailLogin(email string) (*domain.Member, error) {
	member := new(domain.Member)

	emailEncrypted, _ := helpers.Encrypt(email, config.PassPhrase)

	q := fmt.Sprintf(`select "user_id", "first_name", "last_name", "email", "phone", "password", CAST(COALESCE(NULLIF("email_encrypted", ''), '') as VARCHAR)
				from members_aes where email_encrypted = '%s'`, emailEncrypted)

	err := r.db.QueryRow(q).Scan(&member.ID, &member.FirstName, &member.LastName, &member.Email, &member.Phone, &member.Password, &member.EmailEncrypted)

	member.EmailEncrypted, _ = helpers.Decrypt(member.EmailEncrypted, config.PassPhrase)
	return member, err
}
