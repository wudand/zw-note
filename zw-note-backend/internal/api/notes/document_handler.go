package notes

import (
	"strconv"

	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/middleware"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/service"
	"zw-note-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const timeLayout = "2006-01-02T15:04:05Z07:00"

// getNotesUserID reads the user id injected by middleware.NotesAuth.
func getNotesUserID(c *gin.Context) (uint64, bool) {
	v, ok := c.Get(middleware.CtxKeyUserID)
	if !ok {
		return 0, false
	}
	switch id := v.(type) {
	case uint64:
		return id, true
	case float64:
		return uint64(id), true
	default:
		return 0, false
	}
}

// DocumentHandler groups notes document HTTP handlers.
type DocumentHandler struct {
	svc service.DocumentService
	log *zap.Logger
}

func NewDocumentHandler(svc service.DocumentService, log *zap.Logger) *DocumentHandler {
	return &DocumentHandler{svc: svc, log: log}
}

// List godoc
// @Summary  List documents for current user
// @Tags     notes-documents
// @Produce  json
// @Success  200 {object} utils.Response{data=dto.DocumentListResponse}
// @Router   /api/notes/v1/documents [get]
func (h *DocumentHandler) List(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	list, err := h.svc.List(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("list documents", zap.Error(err))
		utils.ServerError(c)
		return
	}

	resp := make([]*dto.DocumentResponse, 0, len(list))
	for _, d := range list {
		resp = append(resp, toDocumentResponse(d))
	}
	utils.Success(c, dto.DocumentListResponse{List: resp})
}

// ListTrash godoc
// @Summary  List soft-deleted documents (trash) for current user
// @Tags     notes-documents
// @Produce  json
// @Success  200 {object} utils.Response{data=dto.DocumentListResponse}
// @Router   /api/notes/v1/documents/trash [get]
func (h *DocumentHandler) ListTrash(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	list, err := h.svc.ListTrash(c.Request.Context(), userID)
	if err != nil {
		h.log.Error("list trashed documents", zap.Error(err))
		utils.ServerError(c)
		return
	}

	resp := make([]*dto.DocumentResponse, 0, len(list))
	for _, d := range list {
		resp = append(resp, toDocumentResponse(d))
	}
	utils.Success(c, dto.DocumentListResponse{List: resp})
}

// Create godoc
// @Summary  Create a new document
// @Tags     notes-documents
// @Accept   json
// @Produce  json
// @Param    body body dto.CreateDocumentRequest true "Create payload"
// @Success  201  {object} utils.Response{data=dto.DocumentResponse}
// @Router   /api/notes/v1/documents [post]
func (h *DocumentHandler) Create(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	var req dto.CreateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	d, err := h.svc.Create(c.Request.Context(), userID, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("create document", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.SuccessCreated(c, toDocumentResponse(d))
}

// GetByID godoc
// @Summary  Get document by ID
// @Tags     notes-documents
// @Produce  json
// @Param    id  path int true "Document ID"
// @Success  200 {object} utils.Response{data=dto.DocumentResponse}
// @Router   /api/notes/v1/documents/{id} [get]
func (h *DocumentHandler) GetByID(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid document id")
		return
	}

	d, err := h.svc.GetByID(c.Request.Context(), id, userID)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get document", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toDocumentResponse(d))
}

// Update godoc
// @Summary  Update a document
// @Tags     notes-documents
// @Accept   json
// @Produce  json
// @Param    id   path int true "Document ID"
// @Param    body body dto.UpdateDocumentRequest true "Update payload"
// @Success  200  {object} utils.Response{data=dto.DocumentResponse}
// @Router   /api/notes/v1/documents/{id} [put]
func (h *DocumentHandler) Update(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid document id")
		return
	}

	var req dto.UpdateDocumentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	d, err := h.svc.Update(c.Request.Context(), id, userID, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("update document", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toDocumentResponse(d))
}

// Delete godoc
// @Summary  Soft-delete a document (moves it to trash, can be restored)
// @Tags     notes-documents
// @Produce  json
// @Param    id  path int true "Document ID"
// @Success  200 {object} utils.Response
// @Router   /api/notes/v1/documents/{id} [delete]
func (h *DocumentHandler) Delete(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid document id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), id, userID); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("delete document", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, nil)
}

// Restore godoc
// @Summary  Restore a soft-deleted document out of the trash
// @Tags     notes-documents
// @Produce  json
// @Param    id  path int true "Document ID"
// @Success  200 {object} utils.Response{data=dto.DocumentResponse}
// @Router   /api/notes/v1/documents/{id}/restore [post]
func (h *DocumentHandler) Restore(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid document id")
		return
	}

	d, err := h.svc.Restore(c.Request.Context(), id, userID)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("restore document", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toDocumentResponse(d))
}

func toDocumentResponse(d *model.Document) *dto.DocumentResponse {
	return &dto.DocumentResponse{
		ID:          strconv.FormatUint(d.ID, 10),
		Title:       d.Title,
		Description: d.Description,
		Author:      d.Author,
		CreatedAt:   d.CreatedAt.Format(timeLayout),
		UpdatedAt:   d.UpdatedAt.Format(timeLayout),
	}
}
