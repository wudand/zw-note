package service

import (
	"context"

	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/repository"
	"zw-note-backend/pkg/utils"

	"go.uber.org/zap"
)

const maxEnabledCarousels = 3

// CarouselService defines business logic for carousel operations.
type CarouselService interface {
	Create(ctx context.Context, req *dto.CreateCarouselRequest) (*model.Carousel, error)
	GetByID(ctx context.Context, id uint64) (*model.Carousel, error)
	List(ctx context.Context, enabledOnly bool) ([]*model.Carousel, error)
	Update(ctx context.Context, id uint64, req *dto.UpdateCarouselRequest) (*model.Carousel, error)
	Delete(ctx context.Context, id uint64) error
}

type carouselService struct {
	repo repository.CarouselRepository
	log  *zap.Logger
}

func NewCarouselService(repo repository.CarouselRepository, log *zap.Logger) CarouselService {
	return &carouselService{repo: repo, log: log}
}

func (s *carouselService) Create(ctx context.Context, req *dto.CreateCarouselRequest) (*model.Carousel, error) {
	c := &model.Carousel{
		ImageURL:  req.ImageURL,
		Title:     req.Title,
		Link:      req.Link,
		SortOrder: req.SortOrder,
		Status:    model.CarouselStatusDisabled,
	}

	if err := s.repo.Create(ctx, c); err != nil {
		return nil, err
	}
	s.log.Info("carousel created", zap.Uint64("id", c.ID))
	return c, nil
}

func (s *carouselService) GetByID(ctx context.Context, id uint64) (*model.Carousel, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, utils.ErrCarouselNotFound
	}
	return c, nil
}

func (s *carouselService) List(ctx context.Context, enabledOnly bool) ([]*model.Carousel, error) {
	return s.repo.List(ctx, enabledOnly)
}

func (s *carouselService) Update(ctx context.Context, id uint64, req *dto.UpdateCarouselRequest) (*model.Carousel, error) {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if c == nil {
		return nil, utils.ErrCarouselNotFound
	}

	if req.ImageURL != nil {
		c.ImageURL = *req.ImageURL
	}
	if req.Title != nil {
		c.Title = *req.Title
	}
	if req.Link != nil {
		c.Link = *req.Link
	}
	if req.SortOrder != nil {
		c.SortOrder = *req.SortOrder
	}
	if req.Status != nil {
		// 启用时检查是否超限
		if *req.Status == model.CarouselStatusEnabled && c.Status == model.CarouselStatusDisabled {
			count, err := s.repo.CountEnabled(ctx)
			if err != nil {
				return nil, err
			}
			if count >= maxEnabledCarousels {
				return nil, utils.ErrCarouselEnabledLimit
			}
		}
		c.Status = *req.Status
	}

	if err := s.repo.Update(ctx, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *carouselService) Delete(ctx context.Context, id uint64) error {
	c, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if c == nil {
		return utils.ErrCarouselNotFound
	}
	return s.repo.Delete(ctx, id)
}
