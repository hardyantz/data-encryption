package domain

import (
	"time"
)

// User data structure
type User struct {
	ID               *string `jsonapi:"primary,memberType" json:"id" query:"id" fieldname:"id"`
	FirstName        *string `jsonapi:"attr,first_name" json:"first_name"  fieldname:"first_name"`
	LastName         *string `jsonapi:"attr,last_name" json:"last_name"  fieldname:"last_name"`
	Password         *string `jsonapi:"attr,password" json:"password"  fieldname:"password"`
	Phone            *string `jsonapi:"attr,phone" json:"phone"  fieldname:"phone"`
	Email            *string `jsonapi:"attr,email" json:"email"  fieldname:"email"`
}

type Parameters struct {
	Limit   string
	OrderBy string
	Sort    string
	Email   string
}

func (u *User) GenID() {
	t := time.Now()
	id := t.Format("20060102150405")
	u.ID = &id
}
