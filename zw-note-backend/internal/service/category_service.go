package service

import (
	"context"

	"go-web-api/internal/dto"
	"go-web-api/internal/model"
	"go-web-api/internal/repository"
	"go-web-api/pkg/utils"

	"go.uber.org/zap"
)

// CategoryService defines business logic for product category operations.
type CategoryService interface {
	Create(ctx context.Context, req *dto.CreateCategoryRequest) (*model.Category, error)
	GetByID(ctx context.Context, id uint64) (*model.Category, error)
	List(ctx context.Context, page, pageSize int, includeDisabled bool) ([]*model.Category, int64, error)
	Update(ctx context.Context, id uint64, req *dto.UpdateCategoryRequest) (*model.Category, error)
	Delete(ctx context.Context, id uint64) error
}

type categoryService struct {
	repo repository.CategoryRepository
	log  *zap.Logger
}

func NewCategoryService(repo repository.CategoryRepository, log *zap.Logger) CategoryService {
	return &categoryService{repo: repo, log: log}
}

func (s *categoryService) Create(ctx context.Context, req *dto.CreateCategoryRequest) (*model.Category, error) {
	cat := &model.Category{
		Name:      req.Name,
		SortOrder: req.SortOrder,
		Status:    1,
	}
	if err := s.repo.Create(ctx, cat); err != nil {
		return nil, err
	}
	s.log.Info("category created", zap.Uint64("id", cat.ID), zap.String("name", cat.Name))
	return cat, nil
}

func (s *categoryService) GetByID(ctx context.Context, id uint64) (*model.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, utils.ErrCategoryNotFound
	}
	return cat, nil
}

func (s *categoryService) List(ctx context.Context, page, pageSize int, includeDisabled bool) ([]*model.Category, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.List(ctx, (page-1)*pageSize, pageSize, includeDisabled)
}

func (s *categoryService) Update(ctx context.Context, id uint64, req *dto.UpdateCategoryRequest) (*model.Category, error) {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, utils.ErrCategoryNotFound
	}

	if req.Name != "" {
		cat.Name = req.Name
	}
	if req.SortOrder != nil {
		cat.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		cat.Status = *req.Status
	}

	if err := s.repo.Update(ctx, cat); err != nil {
		return nil, err
	}
	return cat, nil
}

func (s *categoryService) Delete(ctx context.Context, id uint64) error {
	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if cat == nil {
		return utils.ErrCategoryNotFound
	}
	return s.repo.Delete(ctx, id)
}
