package notes

import (
	"strconv"

	"zw-note-backend/internal/dto"
	"zw-note-backend/internal/model"
	"zw-note-backend/internal/service"
	"zw-note-backend/pkg/utils"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// OutlineHandler groups notes outline HTTP handlers.
type OutlineHandler struct {
	svc service.OutlineService
	log *zap.Logger
}

func NewOutlineHandler(svc service.OutlineService, log *zap.Logger) *OutlineHandler {
	return &OutlineHandler{svc: svc, log: log}
}

// GetTree godoc
// @Summary  Get a document's outline tree
// @Tags     notes-outlines
// @Produce  json
// @Param    id  path int true "Document ID"
// @Success  200 {object} utils.Response{data=[]dto.OutlineNode}
// @Router   /api/notes/v1/documents/{id}/outlines [get]
func (h *OutlineHandler) GetTree(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	documentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid document id")
		return
	}

	tree, err := h.svc.GetTree(c.Request.Context(), documentID, userID)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get outline tree", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, tree)
}

// Create godoc
// @Summary  Create an outline node under a document
// @Tags     notes-outlines
// @Accept   json
// @Produce  json
// @Param    id   path int true "Document ID"
// @Param    body body dto.CreateOutlineRequest true "Create payload"
// @Success  201  {object} utils.Response{data=dto.OutlineNode}
// @Router   /api/notes/v1/documents/{id}/outlines [post]
func (h *OutlineHandler) Create(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	documentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid document id")
		return
	}

	var req dto.CreateOutlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	o, err := h.svc.Create(c.Request.Context(), documentID, userID, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("create outline", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.SuccessCreated(c, toOutlineResponse(o))
}

// Update godoc
// @Summary  Update an outline node
// @Tags     notes-outlines
// @Accept   json
// @Produce  json
// @Param    outlineId  path int true "Outline ID"
// @Param    body       body dto.UpdateOutlineRequest true "Update payload"
// @Success  200        {object} utils.Response{data=dto.OutlineNode}
// @Router   /api/notes/v1/outlines/{outlineId} [put]
func (h *OutlineHandler) Update(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	outlineID, err := strconv.ParseUint(c.Param("outlineId"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid outline id")
		return
	}

	var req dto.UpdateOutlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	o, err := h.svc.Update(c.Request.Context(), outlineID, userID, &req)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("update outline", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toOutlineResponse(o))
}

// Delete godoc
// @Summary  Delete an outline node
// @Tags     notes-outlines
// @Produce  json
// @Param    outlineId  path int true "Outline ID"
// @Success  200        {object} utils.Response
// @Router   /api/notes/v1/outlines/{outlineId} [delete]
func (h *OutlineHandler) Delete(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	outlineID, err := strconv.ParseUint(c.Param("outlineId"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid outline id")
		return
	}

	if err := h.svc.Delete(c.Request.Context(), outlineID, userID); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("delete outline", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, nil)
}

// Reorder godoc
// @Summary  Batch-reorder outline nodes within a document
// @Tags     notes-outlines
// @Accept   json
// @Produce  json
// @Param    id   path int true "Document ID"
// @Param    body body dto.ReorderOutlineRequest true "Reorder payload"
// @Success  200  {object} utils.Response
// @Router   /api/notes/v1/documents/{id}/outlines/reorder [put]
func (h *OutlineHandler) Reorder(c *gin.Context) {
	userID, ok := getNotesUserID(c)
	if !ok {
		utils.Unauthorized(c)
		return
	}

	documentID, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		utils.ParamError(c, "invalid document id")
		return
	}

	var req dto.ReorderOutlineRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	if err := h.svc.Reorder(c.Request.Context(), documentID, userID, req.Items); err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("reorder outlines", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, nil)
}

func toOutlineResponse(o *model.DocumentOutline) *dto.OutlineNode {
	node := &dto.OutlineNode{
		ID:    strconv.FormatUint(o.ID, 10),
		Title: o.Title,
	}
	if o.ParentID != nil {
		parentIDStr := strconv.FormatUint(*o.ParentID, 10)
		node.ParentID = &parentIDStr
	}
	return node
}
