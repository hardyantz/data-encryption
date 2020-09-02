package repository

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/hardyantz/data-encryption/config"
	"github.com/hardyantz/data-encryption/helpers"
	"github.com/hardyantz/data-encryption/pkg/member/domain"
)

type MemberRepository interface {
	Create(member *domain.Member) error
	GetAll(params domain.Parameters) ([]domain.Member, error)
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

var cacheExpired = 1 * time.Hour

func NewMemberRepository(db *sql.DB, conf *config.Config, redis *config.Redis) MemberRepository {
	return &memberRepository{db, conf, redis}
}

func (r *memberRepository) Create(member *domain.Member) error {

	member.GenID()
	query := fmt.Sprintf(`INSERT INTO members(
	id, first_name, last_name, email, phone, password, email_encrypted, email_encrypted_text) 
	VALUES ('%s', '%s', '%s', '%s', '%s', '%s', pgp_pub_encrypt ('%s',dearmor ('%s')), pgp_pub_encrypt ('%s',dearmor ('%s')))`,
		member.ID, member.FirstName, member.LastName, member.Email, member.Phone, member.Password, member.Email, r.conf.PublicKey(), member.Email, r.conf.PublicKey())
	r.db.QueryRow(query)

	getMember, _ := r.GetOne(member.ID)
	member.EmailEncryptText = getMember.EmailEncryptText
	member.EmailEncrypt = getMember.EmailEncrypt
	member.EmailDecrypted = getMember.EmailDecrypted
	member.PhoneEncrypt = getMember.PhoneEncrypt

	// set by ID
	r.redis.Set(member.ID, getMember, cacheExpired)
	r.redis.Set(member.Email, getMember, 0)

	return nil
}

func (r *memberRepository) GetAll(params domain.Parameters) ([]domain.Member, error) {
	var (
		members []domain.Member
		err     error
		sort    string
	)

	keyRedis, err := json.Marshal(params)

	getRedis := r.redis.Get(string(keyRedis))
	if getRedis.Result != nil {
		err := json.Unmarshal(getRedis.Result.([]byte), &members)
		if err == nil {
			return members, nil
		}
	}

	q := fmt.Sprintf(`select "id", "first_name", "last_name", "email", "phone", "password", "email_encrypted", "email_encrypted_text", pgp_pub_decrypt ("email_encrypted",dearmor ('%s')) as email_decrypted  from members`, r.conf.PrivateKey())

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
		if err := rows.Scan(&member.ID, &member.FirstName, &member.LastName, &member.Email, &member.Phone, &member.Password, &member.EmailEncrypt, &member.EmailEncryptText, &member.EmailDecrypted); err != nil {
			return nil, err
		}
		members = append(members, member)
	}

	_ = r.redis.Set(string(keyRedis), members, cacheExpired)

	return members, err
}

func (r *memberRepository) GetOne(id string) (*domain.Member, error) {
	member := new(domain.Member)

	getRedis := r.redis.Get(id)
	if getRedis.Result != nil {
		err := json.Unmarshal(getRedis.Result.([]byte), &member)
		if err == nil {
			return member, nil
		}
	}

	q := fmt.Sprintf(`select "id", "first_name", "last_name", "email", "phone", "password", "email_encrypted", "email_encrypted_text", pgp_pub_decrypt ("email_encrypted",dearmor ('%s')) as email_decrypted from members where id = '%s'`, r.conf.PrivateKey(), id)

	err := r.db.QueryRow(q).Scan(&member.ID, &member.FirstName, &member.LastName, &member.Email, &member.Phone, &member.Password, &member.EmailEncrypt, &member.EmailEncryptText, &member.EmailDecrypted)

	if err != nil {
		return member, err
	}

	_ = r.redis.Set(id, member, cacheExpired)

	return member, err
}

func (r *memberRepository) GetOneEmail(email string) (*domain.Member, error) {
	member := new(domain.Member)

	getRedis := r.redis.Get(email)
	if getRedis.Result != nil {
		err := json.Unmarshal(getRedis.Result.([]byte), &member)
		if err == nil {
			return member, nil
		}
	}

	q := fmt.Sprintf(`select "id", "first_name", "last_name", "email", "phone", "password", "email_encrypted", "email_encrypted_text", pgp_pub_decrypt ("email_encrypted",dearmor ('%s')) as email_decrypted 
				from members where pgp_pub_decrypt ("email_encrypted",dearmor ('%s')) = '%s'`, r.conf.PrivateKey(), r.conf.PrivateKey(), email)

	err := r.db.QueryRow(q).Scan(&member.ID, &member.FirstName, &member.LastName, &member.Email, &member.Phone, &member.Password, &member.EmailEncrypt, &member.EmailEncryptText, &member.EmailDecrypted)

	_ = r.redis.Set(email, member, cacheExpired)

	return member, err
}

func (r *memberRepository) Update(id string, member *domain.Member) error {
	query := fmt.Sprintf(`UPDATE members SET
	"first_name" = '%s', 
	"last_name" = '%s', 
	"phone" = '%s', 
	"password" = '%s'
	WHERE "id" = '%s'`,
		member.FirstName, member.LastName, member.Phone, member.Password, id)

	res, err := r.db.Exec(query)
	if err != nil {
		return err
	}

	if _, err := res.RowsAffected(); err != nil {
		return err
	}

	_ = r.redis.Del(id)

	getMember, _ := r.GetOne(id)
	member.EmailEncryptText = getMember.EmailEncryptText
	member.EmailEncrypt = getMember.EmailEncrypt
	member.EmailDecrypted = getMember.EmailDecrypted
	member.PhoneEncrypt = getMember.PhoneEncrypt
	member.Email = getMember.Email

	_ = r.redis.Set(getMember.Email, getMember, cacheExpired)

	return nil
}

func (r *memberRepository) Delete(id string) error {
	member, err := r.GetOne(id)
	if err != nil {
		return err
	}

	r.db.QueryRow(`DELETE FROM members WHERE ID = '%s'`, id)

	_ = r.redis.Del(id)
	_ = r.redis.Del(member.Email)

	return nil
}

func (r *memberRepository) GetEmailLogin(email string) (*domain.Member, error) {
	member := new(domain.Member)

	getRedis := r.redis.Get(email)
	if getRedis.Result != nil {
		err := json.Unmarshal(getRedis.Result.([]byte), &member)
		if err == nil {
			return member, nil
		}
	}

	hashPassword := helpers.HashAndSalt([]byte(email))

	q := fmt.Sprintf(`select "id", "first_name", "last_name", "email", "phone", "password", "email_encrypted", "email_encrypted_text", pgp_pub_decrypt ("email_encrypted",dearmor ('%s')) as email_decrypted 
				from members where hash_email = '%s'`, r.conf.PrivateKey(), hashPassword)

	err := r.db.QueryRow(q).Scan(&member.ID, &member.FirstName, &member.LastName, &member.Email, &member.Phone, &member.Password, &member.EmailEncrypt, &member.EmailEncryptText, &member.EmailDecrypted)

	_ = r.redis.Set(email, member, 0)

	return member, err
}
