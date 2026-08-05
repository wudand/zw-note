package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"go-web-api/internal/dto"
	"go-web-api/internal/model"
	"go-web-api/internal/repository"
	"go-web-api/pkg/database"
	"go-web-api/pkg/utils"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const (
	redemptionCodeLength = 12
	redemptionCodeChars  = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"
)

// RedemptionCodeService defines business logic for redemption codes.
type RedemptionCodeService interface {
	// Admin
	CreateBatch(ctx context.Context, req *dto.CreateRedemptionCodesRequest) ([]*model.RedemptionCode, error)
	GetByID(ctx context.Context, id uint64) (*model.RedemptionCode, []uint64, error)
	List(ctx context.Context, page, pageSize int, status *int8) ([]*model.RedemptionCode, map[uint64][]uint64, int64, error)
	Update(ctx context.Context, id uint64, req *dto.UpdateRedemptionCodeRequest) (*model.RedemptionCode, error)
	Delete(ctx context.Context, id uint64) error

	// MP
	Validate(ctx context.Context, mpUserID uint64, code string) (*model.RedemptionCode, []*model.Product, error)
	SelectProduct(ctx context.Context, mpUserID uint64, redemptionCodeID, productID uint64) error
}

type redemptionCodeService struct {
	txMgr        database.TransactionManager
	rcRepo       repository.RedemptionCodeRepository
	rcpRepo      repository.RedemptionCodeProductRepository
	productRepo  repository.ProductRepository
	cooldownDays int
	log          *zap.Logger
}

func NewRedemptionCodeService(
	txMgr database.TransactionManager,
	rcRepo repository.RedemptionCodeRepository,
	rcpRepo repository.RedemptionCodeProductRepository,
	productRepo repository.ProductRepository,
	cooldownDays int,
	log *zap.Logger,
) RedemptionCodeService {
	if cooldownDays <= 0 {
		cooldownDays = 30
	}
	return &redemptionCodeService{
		txMgr:        txMgr,
		rcRepo:       rcRepo,
		rcpRepo:      rcpRepo,
		productRepo:  productRepo,
		cooldownDays: cooldownDays,
		log:          log,
	}
}

// generateCode 生成 12 位字母+数字兑换码
func generateCode() (string, error) {
	b := make([]byte, redemptionCodeLength)
	for i := range b {
		n, err := rand.Int(rand.Reader, big.NewInt(int64(len(redemptionCodeChars))))
		if err != nil {
			return "", err
		}
		b[i] = redemptionCodeChars[n.Int64()]
	}
	return string(b), nil
}

func (s *redemptionCodeService) CreateBatch(ctx context.Context, req *dto.CreateRedemptionCodesRequest) ([]*model.RedemptionCode, error) {
	// 校验产品存在且 enabled
	for _, pid := range req.ProductIDs {
		p, err := s.productRepo.GetByID(ctx, pid)
		if err != nil {
			return nil, err
		}
		if p == nil {
			return nil, utils.ErrProductNotFound
		}
		if p.Status != 1 {
			return nil, utils.NewAppError(utils.CodeProductNotFound, "product is disabled")
		}
	}

	seen := make(map[string]bool)
	codes := make([]*model.RedemptionCode, 0, req.Count)
	for len(codes) < req.Count {
		code, err := generateCode()
		if err != nil {
			return nil, err
		}
		if seen[code] {
			continue
		}
		// 检查是否已存在
		existing, err := s.rcRepo.GetByCode(ctx, code)
		if err != nil {
			return nil, err
		}
		if existing != nil {
			continue
		}
		seen[code] = true
		codes = append(codes, &model.RedemptionCode{Code: code, Status: model.RedemptionCodeStatusUnused})
	}

	// 跨 Repo 原子操作：兑换码 + 产品绑定
	if err := s.txMgr.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		if err := s.rcRepo.CreateBatch(ctx, tx, codes); err != nil {
			return err
		}
		for _, rc := range codes {
			if err := s.rcpRepo.CreateBatch(ctx, tx, rc.ID, req.ProductIDs); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, fmt.Errorf("create redemption codes: %w", err)
	}

	s.log.Info("redemption codes created", zap.Int("count", len(codes)), zap.Uint64s("product_ids", req.ProductIDs))
	return codes, nil
}

func (s *redemptionCodeService) GetByID(ctx context.Context, id uint64) (*model.RedemptionCode, []uint64, error) {
	rc, err := s.rcRepo.GetByID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	if rc == nil {
		return nil, nil, utils.ErrRedemptionCodeNotFound
	}
	pids, err := s.rcpRepo.GetProductIDsByCodeID(ctx, id)
	if err != nil {
		return nil, nil, err
	}
	return rc, pids, nil
}

func (s *redemptionCodeService) List(ctx context.Context, page, pageSize int, status *int8) ([]*model.RedemptionCode, map[uint64][]uint64, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	list, total, err := s.rcRepo.List(ctx, (page-1)*pageSize, pageSize, status)
	if err != nil {
		return nil, nil, 0, err
	}
	productMap := make(map[uint64][]uint64, len(list))
	for _, rc := range list {
		pids, _ := s.rcpRepo.GetProductIDsByCodeID(ctx, rc.ID)
		productMap[rc.ID] = pids
	}
	return list, productMap, total, nil
}

func (s *redemptionCodeService) Update(ctx context.Context, id uint64, req *dto.UpdateRedemptionCodeRequest) (*model.RedemptionCode, error) {
	rc, err := s.rcRepo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if rc == nil {
		return nil, utils.ErrRedemptionCodeNotFound
	}
	if rc.Status == model.RedemptionCodeStatusUsed {
		return nil, utils.ErrRedemptionCodeAlreadyUsed
	}

	for _, pid := range req.ProductIDs {
		p, err := s.productRepo.GetByID(ctx, pid)
		if err != nil || p == nil || p.Status != 1 {
			return nil, utils.ErrProductNotFound
		}
	}

	if err := s.txMgr.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		if err := s.rcpRepo.DeleteByCodeID(ctx, tx, id); err != nil {
			return err
		}
		return s.rcpRepo.CreateBatch(ctx, tx, id, req.ProductIDs)
	}); err != nil {
		return nil, fmt.Errorf("update redemption code: %w", err)
	}
	return rc, nil
}

func (s *redemptionCodeService) Delete(ctx context.Context, id uint64) error {
	rc, err := s.rcRepo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if rc == nil {
		return utils.ErrRedemptionCodeNotFound
	}
	if rc.Status == model.RedemptionCodeStatusUsed {
		return utils.ErrRedemptionCodeAlreadyUsed
	}
	return s.rcRepo.Delete(ctx, id)
}

// Validate 校验兑换码并立即锁定给当前用户（抢机制）
func (s *redemptionCodeService) Validate(ctx context.Context, mpUserID uint64, code string) (*model.RedemptionCode, []*model.Product, error) {
	rc, err := s.rcRepo.GetByCode(ctx, code)
	if err != nil {
		return nil, nil, err
	}
	if rc == nil {
		return nil, nil, utils.ErrRedemptionCodeNotFound
	}
	if rc.Status == model.RedemptionCodeStatusUsed {
		// 已被他人使用
		if rc.MPUserID != nil && *rc.MPUserID != mpUserID {
			return nil, nil, utils.ErrRedemptionCodeAlreadyUsed
		}
		// 同一用户已锁定，返回产品列表（幂等）
		pids, err := s.rcpRepo.GetProductIDsByCodeID(ctx, rc.ID)
		if err != nil {
			return nil, nil, err
		}
		products, err := s.loadProducts(ctx, pids)
		if err != nil {
			return nil, nil, err
		}
		return rc, products, nil
	}

	// 防薅：用户最近 30 天内是否使用过任意兑换码
	usedRecently, err := s.rcRepo.HasUserUsedRecently(ctx, mpUserID, s.cooldownDays)
	if err != nil {
		return nil, nil, err
	}
	if usedRecently {
		return nil, nil, utils.ErrRedemptionCodeCooldown
	}

	// 原子抢占：仅当 status=0 时更新，避免并发下多人同时抢到同一兑换码
	claimed, err := s.rcRepo.ClaimIfUnused(ctx, code, mpUserID)
	if err != nil {
		return nil, nil, err
	}
	if !claimed {
		return nil, nil, utils.ErrRedemptionCodeAlreadyUsed
	}

	// 抢占成功，重新获取最新数据
	rc, err = s.rcRepo.GetByCode(ctx, code)
	if err != nil || rc == nil {
		return nil, nil, err
	}

	pids, err := s.rcpRepo.GetProductIDsByCodeID(ctx, rc.ID)
	if err != nil {
		return nil, nil, err
	}
	products, err := s.loadProducts(ctx, pids)
	if err != nil {
		return nil, nil, err
	}

	s.log.Info("redemption code validated and locked", zap.Uint64("id", rc.ID), zap.Uint64("mp_user_id", mpUserID))
	return rc, products, nil
}

func (s *redemptionCodeService) loadProducts(ctx context.Context, pids []uint64) ([]*model.Product, error) {
	products := make([]*model.Product, 0, len(pids))
	for _, pid := range pids {
		p, err := s.productRepo.GetByID(ctx, pid)
		if err != nil {
			return nil, err
		}
		if p != nil && p.Status == 1 {
			products = append(products, p)
		}
	}
	return products, nil
}

// SelectProduct 选择产品（只能调用一次，选定后不可更改）
func (s *redemptionCodeService) SelectProduct(ctx context.Context, mpUserID uint64, redemptionCodeID, productID uint64) error {
	rc, err := s.rcRepo.GetByID(ctx, redemptionCodeID)
	if err != nil {
		return err
	}
	if rc == nil {
		return utils.ErrRedemptionCodeNotFound
	}
	if rc.MPUserID == nil || *rc.MPUserID != mpUserID {
		return utils.ErrRedemptionCodeNotYours
	}
	if rc.UsedProductID != nil {
		return utils.ErrRedemptionCodeProductAlreadySelected
	}

	ok, err := s.rcpRepo.HasProduct(ctx, redemptionCodeID, productID)
	if err != nil {
		return err
	}
	if !ok {
		return utils.ErrProductNotFound
	}

	p, err := s.productRepo.GetByID(ctx, productID)
	if err != nil || p == nil || p.Status != 1 {
		return utils.ErrProductNotFound
	}

	rc.UsedProductID = &productID
	if err := s.rcRepo.Update(ctx, rc); err != nil {
		return err
	}

	s.log.Info("redemption code product selected", zap.Uint64("id", redemptionCodeID), zap.Uint64("product_id", productID))
	return nil
}
