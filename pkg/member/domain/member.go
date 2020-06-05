package domain

import (
	"time"
)

// Member data structure
type Member struct {
	ID               string `jsonapi:"primary,memberType" json:"id" query:"id" fieldname:"id"`
	FirstName        string `jsonapi:"attr,first_name" json:"first_name"  fieldname:"first_name"`
	LastName         string `jsonapi:"attr,last_name" json:"last_name"  fieldname:"last_name"`
	Password         string `jsonapi:"attr,password" json:"password"  fieldname:"password"`
	Phone            string `jsonapi:"attr,phone" json:"phone"  fieldname:"phone"`
	Email            string `jsonapi:"attr,email" json:"email"  fieldname:"email"`
	EmailEncrypt     string `jsonapi:"attr,email_encrypt" json:"email_encrypt"  fieldname:"email_encrypt"`
	PhoneEncrypt     string `jsonapi:"attr,phone_encrypt" json:"phone_encrypt"  fieldname:"phone_encrypt"`
	EmailEncryptText string `jsonapi:"attr,email_encrypted_text" json:"email_encrypted_text"  fieldname:"email_encrypted_text"`
	EmailDecrypted   string `jsonapi:"attr,email_decrypted" json:"email_decrypted"  fieldname:"email_decrypted"`
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
