package notes

import (
	"zw-note-backend/internal/middleware"

	"github.com/gin-gonic/gin"
)

// RegisterRoutes mounts all /api/notes/v1 routes.
//
// MVP has no real login: every request is authenticated via
// middleware.NotesAuth, which injects a fixed default user id.
//
//	GET    /documents                          – list documents
//	POST   /documents                          – create document
//	GET    /documents/:id                      – get document
//	PUT    /documents/:id                      – update document
//	DELETE /documents/:id                      – delete document
//	GET    /documents/:id/outlines             – get outline tree
//	POST   /documents/:id/outlines             – create outline node
//	PUT    /documents/:id/outlines/reorder     – batch-reorder outline nodes
//
//	PUT    /outlines/:outlineId                – update outline node
//	DELETE /outlines/:outlineId                – delete outline node
//	GET    /outlines/:outlineId/content        – get outline content
//	PUT    /outlines/:outlineId/content        – save outline content
//
//	POST   /upload/image                       – upload an image for note content
func RegisterRoutes(r *gin.Engine, docHandler *DocumentHandler, outlineHandler *OutlineHandler, contentHandler *ContentHandler, uploadHandler *UploadHandler, defaultUserID uint64) {
	v1 := r.Group("/api/notes/v1")
	v1.Use(middleware.NotesAuth(defaultUserID))

	docs := v1.Group("/documents")
	docs.GET("", docHandler.List)
	docs.POST("", docHandler.Create)
	docs.GET("/:id", docHandler.GetByID)
	docs.PUT("/:id", docHandler.Update)
	docs.DELETE("/:id", docHandler.Delete)
	docs.GET("/:id/outlines", outlineHandler.GetTree)
	docs.POST("/:id/outlines", outlineHandler.Create)
	docs.PUT("/:id/outlines/reorder", outlineHandler.Reorder)

	outlines := v1.Group("/outlines")
	outlines.PUT("/:outlineId", outlineHandler.Update)
	outlines.DELETE("/:outlineId", outlineHandler.Delete)
	outlines.GET("/:outlineId/content", contentHandler.Get)
	outlines.PUT("/:outlineId/content", contentHandler.Update)

	v1.POST("/upload/image", uploadHandler.UploadImage)
}
