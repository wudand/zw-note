package service

import (
	"context"
	"time"

	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/repository"
	"zw-note-backend/pkg/utils"

	"go.uber.org/zap"
)

// CouponService defines business logic for coupon operations.
type CouponService interface {
	// Admin
	Create(ctx context.Context, req *dto.CreateCouponRequest) (*model.Coupon, error)
	GetByID(ctx context.Context, id uint64) (*model.Coupon, error)
	List(ctx context.Context, page, pageSize int, couponType *string, status *int8) ([]*model.Coupon, int64, error)
	Update(ctx context.Context, id uint64, req *dto.UpdateCouponRequest) (*model.Coupon, error)
	Cancel(ctx context.Context, id uint64) (*model.Coupon, error)
	Delete(ctx context.Context, id uint64) error

	// MP
	ListClaimable(ctx context.Context) ([]*model.Coupon, error)
	Claim(ctx context.Context, mpUserID, couponID uint64) (*model.UserCoupon, error)
	ListMyCoupons(ctx context.Context, mpUserID uint64, status *string) ([]*model.UserCoupon, error)
}

type couponService struct {
	couponRepo     repository.CouponRepository
	userCouponRepo repository.UserCouponRepository
	log            *zap.Logger
}

func NewCouponService(couponRepo repository.CouponRepository, userCouponRepo repository.UserCouponRepository, log *zap.Logger) CouponService {
	return &couponService{couponRepo: couponRepo, userCouponRepo: userCouponRepo, log: log}
}

func (s *couponService) Create(ctx context.Context, req *dto.CreateCouponRequest) (*model.Coupon, error) {
	stackable := int8(0)
	if req.Type == model.CouponTypeNewUser {
		stackable = 1
	}

	c := &model.Coupon{
		Name:          req.Name,
		Type:          req.Type,
		MinAmount:     req.MinAmount,
		DiscountValue: req.DiscountValue,
		ValidDays:     req.ValidDays,
		Stackable:     stackable,
		Status:        model.CouponStatusActive,
	}
	if err := s.couponRepo.Create(ctx, c); err != nil {
		return nil, err
	}
	s.log.Info("coupon created", zap.Uint64("id", c.ID), zap.String("name", c.Name))
	return c, nil
}

func (s *couponService) GetByID(ctx context.Context, id uint64) (*model.Coupon, error) {
	c, err := s.couponRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, utils.ErrCouponNotFound
	}
	return c, nil
}

func (s *couponService) List(ctx context.Context, page, pageSize int, couponType *string, status *int8) ([]*model.Coupon, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.couponRepo.List(ctx, (page-1)*pageSize, pageSize, couponType, status)
}

func (s *couponService) Update(ctx context.Context, id uint64, req *dto.UpdateCouponRequest) (*model.Coupon, error) {
	c, err := s.couponRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, utils.ErrCouponNotFound
	}

	if req.Name != nil {
		c.Name = *req.Name
	}
	if req.MinAmount != nil {
		c.MinAmount = *req.MinAmount
	}
	if req.DiscountValue != nil {
		c.DiscountValue = *req.DiscountValue
	}
	if req.ValidDays != nil {
		c.ValidDays = *req.ValidDays
	}
	if req.Status != nil {
		c.Status = *req.Status
	}

	if err := s.couponRepo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *couponService) Cancel(ctx context.Context, id uint64) (*model.Coupon, error) {
	c, err := s.couponRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, utils.ErrCouponNotFound
	}
	c.Status = model.CouponStatusCancelled
	if err := s.couponRepo.Update(ctx, c); err != nil {
		return nil, err
	}
	s.log.Info("coupon cancelled", zap.Uint64("id", id))
	return c, nil
}

func (s *couponService) Delete(ctx context.Context, id uint64) error {
	c, err := s.couponRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c == nil {
		return utils.ErrCouponNotFound
	}
	return s.couponRepo.Delete(ctx, id)
}

func (s *couponService) ListClaimable(ctx context.Context) ([]*model.Coupon, error) {
	status := model.CouponStatusActive
	list, _, err := s.couponRepo.List(ctx, 0, 1000, nil, &status)
	return list, err
}

func (s *couponService) Claim(ctx context.Context, mpUserID, couponID uint64) (*model.UserCoupon, error) {
	c, err := s.couponRepo.GetByID(ctx, couponID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, utils.ErrCouponNotFound
	}
	if c.Status == model.CouponStatusCancelled {
		return nil, utils.ErrCouponCancelled
	}

	// 新人券：每账号仅用一次，使用后不可再领
	if c.Type == model.CouponTypeNewUser {
		hasUsed, err := s.userCouponRepo.HasUsedNewUserCoupon(ctx, mpUserID)
		if err != nil {
			return nil, err
		}
		if hasUsed {
			return nil, utils.ErrCouponClaimLimit
		}
		// 检查是否已有未使用的新人券
		cnt, err := s.userCouponRepo.CountUnusedByUserAndCoupon(ctx, mpUserID, couponID)
		if err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, utils.ErrCouponClaimLimit
		}
	} else {
		// 满减券：每人限领1张未使用的，使用后可再领
		cnt, err := s.userCouponRepo.CountUnusedByUserAndCoupon(ctx, mpUserID, couponID)
		if err != nil {
			return nil, err
		}
		if cnt > 0 {
			return nil, utils.ErrCouponClaimLimit
		}
	}

	now := time.Now()
	expiryAt := now.AddDate(0, 0, c.ValidDays)

	uc := &model.UserCoupon{
		MPUserID:      mpUserID,
		CouponID:      c.ID,
		CouponName:    c.Name,
		CouponType:    c.Type,
		MinAmount:     c.MinAmount,
		DiscountValue: c.DiscountValue,
		Stackable:     c.Stackable,
		Status:        model.UserCouponStatusUnused,
		ClaimedAt:     now,
		ExpiryAt:      expiryAt,
	}
	if err := s.userCouponRepo.Create(ctx, uc); err != nil {
		return nil, err
	}
	s.log.Info("coupon claimed", zap.Uint64("user_coupon_id", uc.ID), zap.Uint64("mp_user_id", mpUserID), zap.Uint64("coupon_id", couponID))
	return uc, nil
}

func (s *couponService) ListMyCoupons(ctx context.Context, mpUserID uint64, status *string) ([]*model.UserCoupon, error) {
	return s.userCouponRepo.ListByUser(ctx, mpUserID, status)
}
