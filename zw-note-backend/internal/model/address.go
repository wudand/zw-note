package model

import "time"

// Address maps to the `addresses` table.
type Address struct {
	ID        uint64    `db:"id"`
	MPUserID  uint64    `db:"mp_user_id"`
	Receiver  string    `db:"receiver"`
	Phone     string    `db:"phone"`
	Province  string    `db:"province"`
	City      string    `db:"city"`
	District  string    `db:"district"`
	Detail    string    `db:"detail"`
	Tag       string    `db:"tag"`
	IsDefault int8      `db:"is_default"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

const (
	AddressDefault   int8 = 1
	AddressNotDefault int8 = 0
)
