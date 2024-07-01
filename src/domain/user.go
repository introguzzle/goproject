package domain

import (
	"database/sql"
	"time"
)

type User struct {
	Id        int64     `json:"id"`
	Name      string    `json:"name"`
	Email     string    `json:"email"`
	Password  string    `json:"password"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

func (u *User) Table() string {
	return "users"
}

func (u *User) GetId() int64 {
	return u.Id
}

func (u *User) SetId(id int64) {
	u.Id = id
}

func (u *User) FlushRow(row *sql.Row) error {
	return row.Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
}

func (u *User) FlushRows(rows *sql.Rows) error {
	return rows.Scan(
		&u.Id,
		&u.Name,
		&u.Email,
		&u.Password,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
}

func (u *User) Mappings() map[string]string {
	return map[string]string{
		"id":       "Id",
		"name":     "Name",
		"email":    "Email",
		"password": "Password",
	}
}
