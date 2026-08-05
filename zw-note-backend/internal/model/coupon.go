package model

import "time"

// CouponType defines coupon template types.
const (
	CouponTypeNewUser    = "new_user"    // 新人券：每账号仅用一次，可叠加
	CouponTypeSpendReduce = "spend_reduce" // 满减券：满X减Y，不可叠加
)

// Coupon maps to the `coupons` table (template).
type Coupon struct {
	ID            uint64    `db:"id"`
	Name          string    `db:"name"`
	Type          string    `db:"type"`
	MinAmount     int       `db:"min_amount"`     // cents
	DiscountValue int       `db:"discount_value"` // cents
	ValidDays     int       `db:"valid_days"`
	Stackable     int8      `db:"stackable"`
	Status        int8      `db:"status"`
	CreatedAt     time.Time `db:"created_at"`
	UpdatedAt     time.Time `db:"updated_at"`
}

const (
	CouponStatusActive   int8 = 1
	CouponStatusCancelled int8 = 0
)

// UserCoupon maps to the `user_coupons` table (user's claimed coupon).
type UserCoupon struct {
	ID            uint64     `db:"id"`
	MPUserID      uint64     `db:"mp_user_id"`
	CouponID      uint64     `db:"coupon_id"`
	CouponName    string     `db:"coupon_name"`
	CouponType    string     `db:"coupon_type"`
	MinAmount     int        `db:"min_amount"`
	DiscountValue int        `db:"discount_value"`
	Stackable     int8       `db:"stackable"`
	Status        string     `db:"status"`
	ClaimedAt     time.Time  `db:"claimed_at"`
	ExpiryAt      time.Time  `db:"expiry_at"`
	UsedAt        *time.Time `db:"used_at"`
	CreatedAt     time.Time  `db:"created_at"`
	UpdatedAt     time.Time  `db:"updated_at"`
}

const (
	UserCouponStatusUnused  = "unused"
	UserCouponStatusUsed    = "used"
	UserCouponStatusExpired = "expired"
)
