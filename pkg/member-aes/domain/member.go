package domain

import (
	"time"
)

// Member data structure
type Member struct {
	ID             string `jsonapi:"primary,memberType" json:"id" query:"user_id" fieldname:"user_id"`
	FirstName      string `jsonapi:"attr,first_name" json:"first_name"  fieldname:"first_name"`
	LastName       string `jsonapi:"attr,last_name" json:"last_name"  fieldname:"last_name"`
	Password       string `jsonapi:"attr,password" json:"password"  fieldname:"password"`
	Phone          string `jsonapi:"attr,phone" json:"phone"  fieldname:"phone"`
	Email          string `jsonapi:"attr,email" json:"email"  fieldname:"email"`
	EmailEncrypted string `jsonapi:"attr,email_encrypted" json:"email_encrypted" fieldname:"email_encrypted"`
}

type Parameters struct {
	Limit   string
	OrderBy string
	Sort    string
	Email   string
}

func (m *Member) GenID() {
	t := time.Now()
	m.ID = t.Format("20060102150405")
}
