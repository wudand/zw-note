package service

import (
	"context"
	"fmt"

	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/repository"
	"zw-note-backend/pkg/database"
	"zw-note-backend/pkg/utils"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

const maxAddressesPerUser = 10

// AddressService defines business logic for address operations.
type AddressService interface {
	Create(ctx context.Context, mpUserID uint64, req *dto.CreateAddressRequest) (*model.Address, error)
	GetByID(ctx context.Context, id, mpUserID uint64) (*model.Address, error)
	ListByUser(ctx context.Context, mpUserID uint64) ([]*model.Address, error)
	Update(ctx context.Context, id, mpUserID uint64, req *dto.UpdateAddressRequest) (*model.Address, error)
	SetDefault(ctx context.Context, id, mpUserID uint64) (*model.Address, error)
	Delete(ctx context.Context, id, mpUserID uint64) error
}

type addressService struct {
	txMgr database.TransactionManager
	repo  repository.AddressRepository
	log   *zap.Logger
}

func NewAddressService(txMgr database.TransactionManager, repo repository.AddressRepository, log *zap.Logger) AddressService {
	return &addressService{txMgr: txMgr, repo: repo, log: log}
}

func (s *addressService) Create(ctx context.Context, mpUserID uint64, req *dto.CreateAddressRequest) (*model.Address, error) {
	count, err := s.repo.CountByUser(ctx, mpUserID)
	if err != nil {
		return nil, err
	}
	if count >= maxAddressesPerUser {
		return nil, utils.ErrAddressLimitReached
	}

	isDefault := model.AddressNotDefault
	if count == 0 {
		isDefault = model.AddressDefault
	}

	a := &model.Address{
		MPUserID:  mpUserID,
		Receiver:  req.Receiver,
		Phone:     req.Phone,
		Province:  req.Province,
		City:      req.City,
		District:  req.District,
		Detail:    req.Detail,
		Tag:       req.Tag,
		IsDefault: isDefault,
	}
	if err := s.repo.Create(ctx, a); err != nil {
		return nil, err
	}
	s.log.Info("address created", zap.Uint64("id", a.ID), zap.Uint64("mp_user_id", mpUserID))
	return a, nil
}

func (s *addressService) GetByID(ctx context.Context, id, mpUserID uint64) (*model.Address, error) {
	a, err := s.repo.GetByIDAndUser(ctx, id, mpUserID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, utils.ErrAddressNotFound
	}
	return a, nil
}

func (s *addressService) ListByUser(ctx context.Context, mpUserID uint64) ([]*model.Address, error) {
	return s.repo.ListByUser(ctx, mpUserID)
}

func (s *addressService) Update(ctx context.Context, id, mpUserID uint64, req *dto.UpdateAddressRequest) (*model.Address, error) {
	a, err := s.repo.GetByIDAndUser(ctx, id, mpUserID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, utils.ErrAddressNotFound
	}

	if req.Receiver != nil {
		a.Receiver = *req.Receiver
	}
	if req.Phone != nil {
		a.Phone = *req.Phone
	}
	if req.Province != nil {
		a.Province = *req.Province
	}
	if req.City != nil {
		a.City = *req.City
	}
	if req.District != nil {
		a.District = *req.District
	}
	if req.Detail != nil {
		a.Detail = *req.Detail
	}
	if req.Tag != nil {
		a.Tag = *req.Tag
	}

	if err := s.repo.Update(ctx, nil, a); err != nil {
		return nil, err
	}
	return a, nil
}

func (s *addressService) SetDefault(ctx context.Context, id, mpUserID uint64) (*model.Address, error) {
	a, err := s.repo.GetByIDAndUser(ctx, id, mpUserID)
	if err != nil {
		return nil, err
	}
	if a == nil {
		return nil, utils.ErrAddressNotFound
	}
	if a.IsDefault == model.AddressDefault {
		return a, nil
	}

	if err := s.txMgr.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		if err := s.repo.ClearDefaultByUser(ctx, tx, mpUserID); err != nil {
			return err
		}
		a.IsDefault = model.AddressDefault
		return s.repo.Update(ctx, tx, a)
	}); err != nil {
		return nil, fmt.Errorf("set default address: %w", err)
	}
	return a, nil
}

func (s *addressService) Delete(ctx context.Context, id, mpUserID uint64) error {
	a, err := s.repo.GetByIDAndUser(ctx, id, mpUserID)
	if err != nil {
		return err
	}
	if a == nil {
		return utils.ErrAddressNotFound
	}

	if err := s.txMgr.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		if err := s.repo.Delete(ctx, tx, id); err != nil {
			return err
		}
		// If deleted was default, set first remaining as default
		if a.IsDefault == model.AddressDefault {
			list, err := s.repo.ListByUser(ctx, mpUserID)
			if err != nil {
				return err
			}
			if len(list) > 0 {
				first := list[0]
				if err := s.repo.ClearDefaultByUser(ctx, tx, mpUserID); err != nil {
					return err
				}
				first.IsDefault = model.AddressDefault
				return s.repo.Update(ctx, tx, first)
			}
		}
		return nil
	}); err != nil {
		return fmt.Errorf("delete address: %w", err)
	}
	s.log.Info("address deleted", zap.Uint64("id", id), zap.Uint64("mp_user_id", mpUserID))
	return nil
}
