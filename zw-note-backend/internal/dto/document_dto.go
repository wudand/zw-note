package dto

// CreateDocumentRequest is the payload for creating a document.
type CreateDocumentRequest struct {
	Title       string `json:"title"       binding:"required,min=1,max=50"`
	Description string `json:"description" binding:"omitempty,max=200"`
	Author      string `json:"author"      binding:"required,min=1,max=20"`
}

// UpdateDocumentRequest is the payload for updating a document's metadata.
type UpdateDocumentRequest struct {
	Title       *string `json:"title"       binding:"omitempty,min=1,max=50"`
	Description *string `json:"description" binding:"omitempty,max=200"`
	Author      *string `json:"author"      binding:"omitempty,min=1,max=20"`
}

// DocumentResponse is the public-facing representation of a document.
// IDs are serialized as strings to align with the frontend's `id: string` convention.
type DocumentResponse struct {
	ID          string `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Author      string `json:"author"`
	CreatedAt   string `json:"created_at"`
	UpdatedAt   string `json:"updated_at"`
}

// DocumentListResponse wraps a list of documents.
type DocumentListResponse struct {
	List []*DocumentResponse `json:"list"`
}

// CreateOutlineRequest is the payload for creating an outline node under a document.
type CreateOutlineRequest struct {
	Title     string  `json:"title" binding:"required,min=1,max=100"`
	ParentID  *uint64 `json:"parent_id"`
	SortOrder *int    `json:"sort_order"`
}

// UpdateOutlineRequest is the payload for updating an outline node.
// ClearParent, when true, moves the node to the root level regardless of ParentID.
type UpdateOutlineRequest struct {
	Title       *string `json:"title" binding:"omitempty,min=1,max=100"`
	ParentID    *uint64 `json:"parent_id"`
	ClearParent bool    `json:"clear_parent"`
	SortOrder   *int    `json:"sort_order"`
}

// ReorderOutlineItem describes the target parent/order for a single outline node.
type ReorderOutlineItem struct {
	ID        uint64  `json:"id" binding:"required"`
	ParentID  *uint64 `json:"parent_id"`
	SortOrder int     `json:"sort_order"`
}

// ReorderOutlineRequest is the payload for batch-reordering outline nodes.
type ReorderOutlineRequest struct {
	Items []ReorderOutlineItem `json:"items" binding:"required,dive"`
}

// OutlineNode is the tree-shaped response for a document's outline.
// IDs are serialized as strings to align with the frontend's `OutlineItem` type.
type OutlineNode struct {
	ID       string         `json:"id"`
	Title    string         `json:"title"`
	ParentID *string        `json:"parent_id,omitempty"`
	Children []*OutlineNode `json:"children,omitempty"`
}

// UpdateOutlineContentRequest is the payload for saving an outline node's Markdown content.
type UpdateOutlineContentRequest struct {
	Content string `json:"content"`
}

// OutlineContentResponse is the public-facing representation of an outline node's content.
type OutlineContentResponse struct {
	OutlineID string `json:"outline_id"`
	Content   string `json:"content"`
	UpdatedAt string `json:"updated_at"`
}
