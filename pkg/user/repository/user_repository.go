package repository

import (
	"database/sql"
	"fmt"
	"strings"

	"github.com/hardyantz/data-encryption/config"
	"github.com/hardyantz/data-encryption/helpers"
	"github.com/hardyantz/data-encryption/pkg/user/domain"
)

type UserImplementation interface {
	Create(user *domain.User) error
	GetAll(params domain.Parameters) ([]domain.User, error)
	GetOne(id string) (domain.User, error)
	GetOneEmail(email string) (domain.User, error)
	Update(id string, user *domain.User) error
	Delete(id string) error
}

type UserRepository struct {
	db   *sql.DB
	conf *config.Config
}

func NewUserImplementation(db *sql.DB, conf *config.Config) *UserRepository {
	return &UserRepository{db, conf}
}

func (r *UserRepository) Create(user *domain.User) error {
	user.GenID()
	query := fmt.Sprintf(`INSERT INTO users(
	id, first_name, last_name, email, phone, password) 
	VALUES ('%s', '%s', '%s', '%s', '%s', '%s')`,
		user.ID, user.FirstName, user.LastName, user.Email, user.Phone, user.Password)
	r.db.QueryRow(query)

	return nil
}

func (r *UserRepository) GetAll(params domain.Parameters) ([]domain.User, error) {
	var (
		members []domain.User
		err     error
		sort    string
		limit   string
	)

	q := fmt.Sprintf(`select "id", "first_name", "last_name", "email", "phone", "password" from users`)

	sort = "asc"

	if len(params.OrderBy) > 0 && len(params.Sort) > 0 {
		if helpers.StringInSlice(strings.ToLower(params.Sort), []string{"asc", "desc"}) {
			sort = params.Sort
		}
		q += fmt.Sprintf(` ORDER BY %s %s`, params.OrderBy, sort)
	}

	limit = "10"
	if len(params.Limit) > 0 {
		limit = params.Limit
	}

	q += fmt.Sprintf(" LIMIT %s", limit)

	rows, err := r.db.Query(q)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var user domain.User
		if err := rows.Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password); err != nil {
			return nil, err
		}
		members = append(members, user)
	}

	return members, err
}

func (r *UserRepository) GetOne(id string) (domain.User, error) {
	var user domain.User

	q := fmt.Sprintf(`select "id", "first_name", "last_name", "email", "phone", "password" from users where id = '%s'`, id)

	err := r.db.QueryRow(q).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password)

	return user, err
}

func (r *UserRepository) GetOneEmail(email string) (domain.User, error) {
	var user domain.User

	q := fmt.Sprintf(`select "id", "first_name", "last_name", "email", "phone", "password"
				from users where email = '%s'`, email)

	err := r.db.QueryRow(q).Scan(&user.ID, &user.FirstName, &user.LastName, &user.Email, &user.Phone, &user.Password)

	return user, err
}

func (r *UserRepository) Update(id string, user *domain.User) error {
	query := fmt.Sprintf(`UPDATE users SET
	"first_name" = '%s', 
	"last_name" = '%s', 
	"email" = '%s', 
	"phone" = '%s', 
	"password" = '%s'
	WHERE "id" = '%s'`,
		user.FirstName, user.LastName, user.Email, user.Phone, user.Password, id)
	r.db.QueryRow(query)

	return nil
}

func (r *UserRepository) Delete(id string) error {
	r.db.QueryRow(`DELETE FROM users WHERE ID = '%s'`, id)
	return nil
}
