package service

import (
	"context"
	"strconv"
	"time"

	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/repository"
	"zw-note-backend/pkg/database"
	"zw-note-backend/pkg/utils"

	"github.com/jmoiron/sqlx"
	"go.uber.org/zap"
)

// OutlineService defines business logic for outline nodes and their content.
// It manages both `document_outlines` and `outline_contents` because the two
// are a tightly-coupled 1:1 relationship.
type OutlineService interface {
	GetTree(ctx context.Context, documentID, userID uint64) ([]*dto.OutlineNode, error)
	Create(ctx context.Context, documentID, userID uint64, req *dto.CreateOutlineRequest) (*model.DocumentOutline, error)
	Update(ctx context.Context, outlineID, userID uint64, req *dto.UpdateOutlineRequest) (*model.DocumentOutline, error)
	Delete(ctx context.Context, outlineID, userID uint64) error
	Reorder(ctx context.Context, documentID, userID uint64, items []dto.ReorderOutlineItem) error
	GetContent(ctx context.Context, outlineID, userID uint64) (*model.OutlineContent, error)
	SaveContent(ctx context.Context, outlineID, userID uint64, content string) (*model.OutlineContent, error)
}

type outlineService struct {
	txMgr       database.TransactionManager
	outlineRepo repository.OutlineRepository
	contentRepo repository.OutlineContentRepository
	docRepo     repository.DocumentRepository
	log         *zap.Logger
}

func NewOutlineService(
	txMgr database.TransactionManager,
	outlineRepo repository.OutlineRepository,
	contentRepo repository.OutlineContentRepository,
	docRepo repository.DocumentRepository,
	log *zap.Logger,
) OutlineService {
	return &outlineService{
		txMgr:       txMgr,
		outlineRepo: outlineRepo,
		contentRepo: contentRepo,
		docRepo:     docRepo,
		log:         log,
	}
}

// requireOwnedDocument checks that documentID belongs to userID.
func (s *outlineService) requireOwnedDocument(ctx context.Context, documentID, userID uint64) error {
	doc, err := s.docRepo.GetByIDAndUser(ctx, documentID, userID)
	if err != nil {
		return err
	}
	if doc == nil {
		return utils.ErrDocumentNotFound
	}
	return nil
}

// getOwnedOutline loads the outline node and verifies it belongs to a document
// owned by userID. document_outlines has no user_id column of its own, so
// ownership is always checked transitively through the parent document.
func (s *outlineService) getOwnedOutline(ctx context.Context, outlineID, userID uint64) (*model.DocumentOutline, error) {
	o, err := s.outlineRepo.GetByID(ctx, outlineID)
	if err != nil {
		return nil, err
	}
	if o == nil {
		return nil, utils.ErrOutlineNotFound
	}
	doc, err := s.docRepo.GetByIDAndUser(ctx, o.DocumentID, userID)
	if err != nil {
		return nil, err
	}
	if doc == nil {
		return nil, utils.ErrOutlineNotFound
	}
	return o, nil
}

// GetTree loads all nodes of a document and assembles them into a tree,
// rooted at nodes whose parent_id is nil.
func (s *outlineService) GetTree(ctx context.Context, documentID, userID uint64) ([]*dto.OutlineNode, error) {
	if err := s.requireOwnedDocument(ctx, documentID, userID); err != nil {
		return nil, err
	}

	list, err := s.outlineRepo.ListByDocument(ctx, documentID)
	if err != nil {
		return nil, err
	}

	nodes := make(map[uint64]*dto.OutlineNode, len(list))
	for _, o := range list {
		nodes[o.ID] = &dto.OutlineNode{
			ID:    strconv.FormatUint(o.ID, 10),
			Title: o.Title,
		}
	}

	roots := make([]*dto.OutlineNode, 0)
	for _, o := range list {
		node := nodes[o.ID]
		if o.ParentID == nil {
			roots = append(roots, node)
			continue
		}
		parentIDStr := strconv.FormatUint(*o.ParentID, 10)
		node.ParentID = &parentIDStr
		if parent, ok := nodes[*o.ParentID]; ok {
			parent.Children = append(parent.Children, node)
		} else {
			// Dangling parent reference (shouldn't happen given FK constraints); surface as root.
			roots = append(roots, node)
		}
	}
	return roots, nil
}

// Create adds a new outline node under documentID. If ParentID is set, it
// must reference an existing node within the same document.
func (s *outlineService) Create(ctx context.Context, documentID, userID uint64, req *dto.CreateOutlineRequest) (*model.DocumentOutline, error) {
	if err := s.requireOwnedDocument(ctx, documentID, userID); err != nil {
		return nil, err
	}

	if req.ParentID != nil {
		parent, err := s.outlineRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		if parent == nil || parent.DocumentID != documentID {
			return nil, utils.ErrOutlineInvalidParent
		}
	}

	sortOrder := 0
	if req.SortOrder != nil {
		sortOrder = *req.SortOrder
	}

	o := &model.DocumentOutline{
		DocumentID: documentID,
		ParentID:   req.ParentID,
		Title:      req.Title,
		SortOrder:  sortOrder,
	}

	if err := s.txMgr.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		if err := s.outlineRepo.Create(ctx, tx, o); err != nil {
			return err
		}
		return s.contentRepo.Upsert(ctx, tx, o.ID, "")
	}); err != nil {
		return nil, err
	}

	return o, nil
}

// Update modifies an outline node's title, sort order, and/or parent.
// Moving a node validates that the new parent exists in the same document and
// that the move would not create a cycle (i.e. the new parent is not a
// descendant of the node being moved).
func (s *outlineService) Update(ctx context.Context, outlineID, userID uint64, req *dto.UpdateOutlineRequest) (*model.DocumentOutline, error) {
	o, err := s.getOwnedOutline(ctx, outlineID, userID)
	if err != nil {
		return nil, err
	}

	if req.Title != nil {
		o.Title = *req.Title
	}
	if req.SortOrder != nil {
		o.SortOrder = *req.SortOrder
	}

	switch {
	case req.ClearParent:
		o.ParentID = nil
	case req.ParentID != nil:
		if *req.ParentID == o.ID {
			return nil, utils.ErrOutlineInvalidParent
		}
		newParent, err := s.outlineRepo.GetByID(ctx, *req.ParentID)
		if err != nil {
			return nil, err
		}
		if newParent == nil || newParent.DocumentID != o.DocumentID {
			return nil, utils.ErrOutlineInvalidParent
		}
		isDescendant, err := s.outlineRepo.IsDescendant(ctx, o.DocumentID, o.ID, *req.ParentID)
		if err != nil {
			return nil, err
		}
		if isDescendant {
			return nil, utils.ErrOutlineInvalidParent
		}
		o.ParentID = req.ParentID
	}

	if err := s.outlineRepo.Update(ctx, nil, o); err != nil {
		return nil, err
	}
	return o, nil
}

// Delete removes an outline node; descendant nodes and content are cleaned up
// by the FK's ON DELETE CASCADE.
func (s *outlineService) Delete(ctx context.Context, outlineID, userID uint64) error {
	if _, err := s.getOwnedOutline(ctx, outlineID, userID); err != nil {
		return err
	}
	return s.outlineRepo.Delete(ctx, nil, outlineID)
}

// Reorder batch-updates the parent/sort_order of multiple nodes within one
// document, all-or-nothing. Each move is validated the same way as Update:
// the target parent must belong to the same document and must not be a
// descendant of the node being moved.
func (s *outlineService) Reorder(ctx context.Context, documentID, userID uint64, items []dto.ReorderOutlineItem) error {
	if err := s.requireOwnedDocument(ctx, documentID, userID); err != nil {
		return err
	}

	return s.txMgr.WithTransaction(ctx, func(tx *sqlx.Tx) error {
		for _, item := range items {
			o, err := s.outlineRepo.GetByID(ctx, item.ID)
			if err != nil {
				return err
			}
			if o == nil || o.DocumentID != documentID {
				return utils.ErrOutlineInvalidParent
			}

			if item.ParentID != nil {
				if *item.ParentID == item.ID {
					return utils.ErrOutlineInvalidParent
				}
				parent, err := s.outlineRepo.GetByID(ctx, *item.ParentID)
				if err != nil {
					return err
				}
				if parent == nil || parent.DocumentID != documentID {
					return utils.ErrOutlineInvalidParent
				}
				// Reads run outside tx (read-committed), so this still reflects the
				// pre-reorder tree shape and is unaffected by earlier updates in this
				// same batch.
				isDescendant, err := s.outlineRepo.IsDescendant(ctx, documentID, item.ID, *item.ParentID)
				if err != nil {
					return err
				}
				if isDescendant {
					return utils.ErrOutlineInvalidParent
				}
			}

			o.ParentID = item.ParentID
			o.SortOrder = item.SortOrder
			if err := s.outlineRepo.Update(ctx, tx, o); err != nil {
				return err
			}
		}
		return nil
	})
}

func (s *outlineService) GetContent(ctx context.Context, outlineID, userID uint64) (*model.OutlineContent, error) {
	if _, err := s.getOwnedOutline(ctx, outlineID, userID); err != nil {
		return nil, err
	}

	c, err := s.contentRepo.GetByOutlineID(ctx, outlineID)
	if err != nil {
		return nil, err
	}
	if c == nil {
		// Should not normally happen (content is created alongside the outline),
		// but degrade gracefully instead of erroring.
		c = &model.OutlineContent{OutlineID: outlineID, Content: ""}
	}
	return c, nil
}

func (s *outlineService) SaveContent(ctx context.Context, outlineID, userID uint64, content string) (*model.OutlineContent, error) {
	if _, err := s.getOwnedOutline(ctx, outlineID, userID); err != nil {
		return nil, err
	}

	if err := s.contentRepo.Upsert(ctx, nil, outlineID, content); err != nil {
		return nil, err
	}
	return &model.OutlineContent{OutlineID: outlineID, Content: content, UpdatedAt: time.Now()}, nil
}
