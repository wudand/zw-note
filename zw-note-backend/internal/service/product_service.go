package service

import (
	"context"

	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/repository"
	"zw-note-backend/pkg/types"
	"zw-note-backend/pkg/utils"

	"go.uber.org/zap"
)

// ProductService defines business logic for product operations.
type ProductService interface {
	Create(ctx context.Context, req *dto.CreateProductRequest) (*model.Product, error)
	GetByID(ctx context.Context, id uint64) (*model.Product, error)
	List(ctx context.Context, page, pageSize int, categoryID *uint64, includeDisabled bool) ([]*model.Product, int64, error)
	Update(ctx context.Context, id uint64, req *dto.UpdateProductRequest) (*model.Product, error)
	Delete(ctx context.Context, id uint64) error
}

type productService struct {
	repo         repository.ProductRepository
	categoryRepo repository.CategoryRepository
	log          *zap.Logger
}

func NewProductService(repo repository.ProductRepository, categoryRepo repository.CategoryRepository, log *zap.Logger) ProductService {
	return &productService{repo: repo, categoryRepo: categoryRepo, log: log}
}

func (s *productService) Create(ctx context.Context, req *dto.CreateProductRequest) (*model.Product, error) {
	cat, err := s.categoryRepo.GetByID(ctx, req.CategoryID)
	if err != nil {
		return nil, err
	}
	if cat == nil {
		return nil, utils.ErrCategoryNotFound
	}

	p := &model.Product{
		CategoryID:     req.CategoryID,
		Name:           req.Name,
		Ingredients:    req.Ingredients,
		Grade:          req.Grade,
		Storage:        req.Storage,
		Specification:  req.Specification,
		CoverImage:     req.CoverImage,
		CarouselImages: types.StringSlice(req.CarouselImages),
		Detail:         req.Detail,
		DetailImages:   types.StringSlice(req.DetailImages),
		Status:         1,
	}
	if err := s.repo.Create(ctx, p); err != nil {
		return nil, err
	}
	s.log.Info("product created", zap.Uint64("id", p.ID), zap.String("name", p.Name))
	return p, nil
}

func (s *productService) GetByID(ctx context.Context, id uint64) (*model.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, utils.ErrProductNotFound
	}
	return p, nil
}

func (s *productService) List(ctx context.Context, page, pageSize int, categoryID *uint64, includeDisabled bool) ([]*model.Product, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.List(ctx, (page-1)*pageSize, pageSize, categoryID, includeDisabled)
}

func (s *productService) Update(ctx context.Context, id uint64, req *dto.UpdateProductRequest) (*model.Product, error) {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if p == nil {
		return nil, utils.ErrProductNotFound
	}

	if req.CategoryID != nil {
		cat, err := s.categoryRepo.GetByID(ctx, *req.CategoryID)
		if err != nil {
			return nil, err
		}
		if cat == nil {
			return nil, utils.ErrCategoryNotFound
		}
		p.CategoryID = *req.CategoryID
	}
	if req.Name != nil {
		p.Name = *req.Name
	}
	if req.Ingredients != nil {
		p.Ingredients = *req.Ingredients
	}
	if req.Grade != nil {
		p.Grade = *req.Grade
	}
	if req.Storage != nil {
		p.Storage = *req.Storage
	}
	if req.Specification != nil {
		p.Specification = *req.Specification
	}
	if req.CoverImage != nil {
		p.CoverImage = *req.CoverImage
	}
	if req.CarouselImages != nil {
		p.CarouselImages = types.StringSlice(req.CarouselImages)
	}
	if req.Detail != nil {
		p.Detail = *req.Detail
	}
	if req.DetailImages != nil {
		p.DetailImages = types.StringSlice(req.DetailImages)
	}
	if req.Status != nil {
		p.Status = *req.Status
	}

	if err := s.repo.Update(ctx, p); err != nil {
		return nil, err
	}
	return p, nil
}

func (s *productService) Delete(ctx context.Context, id uint64) error {
	p, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if p == nil {
		return utils.ErrProductNotFound
	}
	return s.repo.Delete(ctx, id)
}
