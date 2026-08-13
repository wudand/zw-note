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

// defaultOutlineTitle is the title given to the root outline node auto-created
// alongside a new document, so the edit page has something to open immediately.
const defaultOutlineTitle = "正文"

// DocumentService defines business logic for document operations.
type DocumentService interface {
	Create(ctx context.Context, userID uint64, req *dto.CreateDocumentRequest) (*model.Document, error)
	GetByID(ctx context.Context, id, userID uint64) (*model.Document, error)
	List(ctx context.Context, userID uint64) ([]*model.Document, error)
	Update(ctx context.Context, id, userID uint64, req *dto.UpdateDocumentRequest) (*model.Document, error)
	// Delete soft-deletes a document (sets deleted_at); it can be undone via Restore.
	Delete(ctx context.Context, id, userID uint64) error
	// ListTrash lists the current user's soft-deleted documents.
	ListTrash(ctx context.Context, userID uint64) ([]*model.Document, error)
	// Restore undoes a soft delete, bringing the document (and its untouched
	// outlines/contents) back into the active list.
	Restore(ctx context.Context, id, userID uint64) (*model.Document, error)
}

type documentService struct {
	txMgr       database.TransactionManager
	docRepo     repository.DocumentRepository
	outlineRepo repository.OutlineRepository
	contentRepo repository.OutlineContentRepository
	log         *zap.Logger
}

func NewDocumentService(
	txMgr database.TransactionManager,
	docRepo repository.DocumentRepository,
	outlineRepo repository.OutlineRepository,
	contentRepo repository.OutlineContentRepository,
	log *zap.Logger,
) DocumentService {
	return &documentService{
		txMgr:       txMgr,
		docRepo:     docRepo,
		outlineRepo: outlineRepo,
		contentRepo: contentRepo,
		log:         log,
	}
}

// Create creates a document together with a default root outline node and its
// empty content, all within a single transaction, so the caller can jump
// straight into editing after creation.
func (s *documentService) Create(ctx context.Context, userID uint64, req *dto.CreateDocumentRequest) (*model.Document, error) {
	d := &model.Document{
		UserID:      userID,
		Title:       req.Title,
		Description: req.Description,
		Author:      req.Author,
	}

	if err := s.txMgr.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		if err := s.docRepo.Create(ctx, tx, d); err != nil {
			return err
		}

		outline := &model.DocumentOutline{
			DocumentID: d.ID,
			ParentID:   nil,
			Title:      defaultOutlineTitle,
			SortOrder:  0,
		}
		if err := s.outlineRepo.Create(ctx, tx, outline); err != nil {
			return err
		}

		return s.contentRepo.Upsert(ctx, tx, outline.ID, "")
	}); err != nil {
		return nil, fmt.Errorf("create document: %w", err)
	}

	s.log.Info("document created", zap.Uint64("id", d.ID), zap.Uint64("user_id", userID))
	return d, nil
}

func (s *documentService) GetByID(ctx context.Context, id, userID uint64) (*model.Document, error) {
	d, err := s.docRepo.GetByIDAndUser(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, utils.ErrDocumentNotFound
	}
	return d, nil
}

func (s *documentService) List(ctx context.Context, userID uint64) ([]*model.Document, error) {
	return s.docRepo.ListByUser(ctx, userID)
}

func (s *documentService) Update(ctx context.Context, id, userID uint64, req *dto.UpdateDocumentRequest) (*model.Document, error) {
	d, err := s.docRepo.GetByIDAndUser(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, utils.ErrDocumentNotFound
	}

	if req.Title != nil {
		d.Title = *req.Title
	}
	if req.Description != nil {
		d.Description = *req.Description
	}
	if req.Author != nil {
		d.Author = *req.Author
	}

	if err := s.docRepo.Update(ctx, nil, d); err != nil {
		return nil, err
	}
	return d, nil
}

func (s *documentService) Delete(ctx context.Context, id, userID uint64) error {
	d, err := s.docRepo.GetByIDAndUser(ctx, id, userID)
	if err != nil {
		return err
	}
	if d == nil {
		return utils.ErrDocumentNotFound
	}

	// Logical delete only: document_outlines / outline_contents are left untouched
	// so Restore can bring everything back exactly as it was.
	if err := s.docRepo.SoftDelete(ctx, nil, id, userID); err != nil {
		return err
	}
	s.log.Info("document soft-deleted", zap.Uint64("id", id), zap.Uint64("user_id", userID))
	return nil
}

func (s *documentService) ListTrash(ctx context.Context, userID uint64) ([]*model.Document, error) {
	return s.docRepo.ListDeletedByUser(ctx, userID)
}

func (s *documentService) Restore(ctx context.Context, id, userID uint64) (*model.Document, error) {
	d, err := s.docRepo.GetDeletedByIDAndUser(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	if d == nil {
		return nil, utils.ErrDocumentNotFound
	}

	if err := s.docRepo.Restore(ctx, nil, id, userID); err != nil {
		return nil, err
	}

	restored, err := s.docRepo.GetByIDAndUser(ctx, id, userID)
	if err != nil {
		return nil, err
	}
	s.log.Info("document restored", zap.Uint64("id", id), zap.Uint64("user_id", userID))
	return restored, nil
}
