package domain

import (
	"database/sql"
	"time"
)

type Product struct {
	Id          int64     `json:"id"`
	Name        string    `json:"name"`
	Description string    `json:"description"`
	Price       float64   `json:"price"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func (p *Product) Table() string {
	return "products"
}

func (p *Product) GetId() int64 {
	return p.Id
}

func (p *Product) SetId(id int64) {
	p.Id = id
}

func (p *Product) FlushRow(row *sql.Row) error {
	return row.Scan(
		&p.Id,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
}

func (p *Product) FlushRows(rows *sql.Rows) error {
	return rows.Scan(
		&p.Id,
		&p.Name,
		&p.Description,
		&p.Price,
		&p.CreatedAt,
		&p.UpdatedAt,
	)
}

func (p *Product) Mappings() map[string]string {
	return map[string]string{
		"id":          "Id",
		"name":        "Name",
		"description": "Description",
		"price":       "Price",
		"created_at":  "CreatedAt",
		"updated_at":  "UpdatedAt",
	}
}
