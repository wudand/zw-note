package service

import (
	"context"
	"fmt"

	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/repository"
	"zw-note-backend/pkg/utils"

	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

//go:generate mockgen -destination=mock/admin_user_service_mock.go -package=mock zw-note-backend/internal/service AdminUserService

// AdminUserService defines business logic for management-console user operations.
type AdminUserService interface {
	Create(ctx context.Context, req *dto.CreateAdminUserRequest) (*model.AdminUser, error)
	Authenticate(ctx context.Context, username, password string) (*model.AdminUser, error)
	GetByID(ctx context.Context, id uint64) (*model.AdminUser, error)
	List(ctx context.Context, page, pageSize int) ([]*model.AdminUser, int64, error)
	Update(ctx context.Context, id uint64, req *dto.UpdateAdminUserRequest) (*model.AdminUser, error)
	Delete(ctx context.Context, id uint64) error
}

type adminUserService struct {
	repo repository.AdminUserRepository
	log  *zap.Logger
}

func NewAdminUserService(repo repository.AdminUserRepository, log *zap.Logger) AdminUserService {
	return &adminUserService{repo: repo, log: log}
}

func (s *adminUserService) Create(ctx context.Context, req *dto.CreateAdminUserRequest) (*model.AdminUser, error) {
	existing, err := s.repo.GetByUsername(ctx, req.Username)
	if err != nil {
		return nil, fmt.Errorf("check username: %w", err)
	}
	if existing != nil {
		return nil, utils.ErrUserExists
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	user := &model.AdminUser{
		Username: req.Username,
		Password: string(hashed),
		Role:     req.Role,
		Status:   1,
	}

	if err := s.repo.Create(ctx, user); err != nil {
		return nil, err
	}

	s.log.Info("admin user created", zap.String("username", user.Username), zap.String("role", user.Role))
	return user, nil
}

func (s *adminUserService) Authenticate(ctx context.Context, username, password string) (*model.AdminUser, error) {
	user, err := s.repo.GetByUsername(ctx, username)
	if err != nil {
		return nil, fmt.Errorf("get admin user: %w", err)
	}
	if user == nil {
		return nil, utils.ErrUserNotFound
	}
	if user.Status == 0 {
		return nil, utils.ErrUserDisabled
	}
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, utils.ErrInvalidPassword
	}
	return user, nil
}

func (s *adminUserService) GetByID(ctx context.Context, id uint64) (*model.AdminUser, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.ErrUserNotFound
	}
	return user, nil
}

func (s *adminUserService) List(ctx context.Context, page, pageSize int) ([]*model.AdminUser, int64, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 10
	}
	return s.repo.List(ctx, (page-1)*pageSize, pageSize)
}

func (s *adminUserService) Update(ctx context.Context, id uint64, req *dto.UpdateAdminUserRequest) (*model.AdminUser, error) {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, utils.ErrUserNotFound
	}

	if req.Username != "" {
		user.Username = req.Username
	}
	if req.Role != "" {
		user.Role = req.Role
	}
	if req.Status != nil {
		user.Status = *req.Status
	}

	if err := s.repo.Update(ctx, user); err != nil {
		return nil, err
	}
	return user, nil
}

func (s *adminUserService) Delete(ctx context.Context, id uint64) error {
	user, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return err
	}
	if user == nil {
		return utils.ErrUserNotFound
	}
	return s.repo.Delete(ctx, id)
}
