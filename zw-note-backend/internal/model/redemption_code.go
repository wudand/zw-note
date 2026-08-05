package model

import "time"

const (
	RedemptionCodeStatusUnused int8 = 0
	RedemptionCodeStatusUsed  int8 = 1
)

// RedemptionCode maps to the `redemption_codes` table.
type RedemptionCode struct {
	ID             uint64     `db:"id"`
	Code           string     `db:"code"`
	Status         int8       `db:"status"`
	MPUserID       *uint64    `db:"mp_user_id"`
	UsedProductID  *uint64    `db:"used_product_id"`
	UsedAt         *time.Time `db:"used_at"`
	CreatedAt      time.Time  `db:"created_at"`
	UpdatedAt      time.Time  `db:"updated_at"`
}

// RedemptionCodeProduct maps to the `redemption_code_products` table.
type RedemptionCodeProduct struct {
	ID               uint64    `db:"id"`
	RedemptionCodeID uint64    `db:"redemption_code_id"`
	ProductID        uint64    `db:"product_id"`
	CreatedAt        time.Time `db:"created_at"`
}
