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

// ContentHandler groups notes outline-content HTTP handlers.
type ContentHandler struct {
	svc service.OutlineService
	log *zap.Logger
}

func NewContentHandler(svc service.OutlineService, log *zap.Logger) *ContentHandler {
	return &ContentHandler{svc: svc, log: log}
}

// Get godoc
// @Summary  Get an outline node's Markdown content
// @Tags     notes-contents
// @Produce  json
// @Param    outlineId  path int true "Outline ID"
// @Success  200        {object} utils.Response{data=dto.OutlineContentResponse}
// @Router   /api/notes/v1/outlines/{outlineId}/content [get]
func (h *ContentHandler) Get(c *gin.Context) {
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

	content, err := h.svc.GetContent(c.Request.Context(), outlineID, userID)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("get outline content", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toContentResponse(content))
}

// Update godoc
// @Summary  Save an outline node's Markdown content
// @Tags     notes-contents
// @Accept   json
// @Produce  json
// @Param    outlineId  path int true "Outline ID"
// @Param    body       body dto.UpdateOutlineContentRequest true "Content payload"
// @Success  200        {object} utils.Response{data=dto.OutlineContentResponse}
// @Router   /api/notes/v1/outlines/{outlineId}/content [put]
func (h *ContentHandler) Update(c *gin.Context) {
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

	var req dto.UpdateOutlineContentRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		utils.ParamError(c, err.Error())
		return
	}

	content, err := h.svc.SaveContent(c.Request.Context(), outlineID, userID, req.Content)
	if err != nil {
		if !utils.HandleAppError(c, err) {
			h.log.Error("save outline content", zap.Error(err))
			utils.ServerError(c)
		}
		return
	}

	utils.Success(c, toContentResponse(content))
}

func toContentResponse(o *model.OutlineContent) *dto.OutlineContentResponse {
	return &dto.OutlineContentResponse{
		OutlineID: strconv.FormatUint(o.OutlineID, 10),
		Content:   o.Content,
		UpdatedAt: o.UpdatedAt.Format(timeLayout),
	}
}
