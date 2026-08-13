-- Migration: 006_add_documents_deleted_at
-- Direction: UP
-- 文档删除改为逻辑删除：新增 deleted_at，NULL 表示未删除，
-- 删除时只回填时间戳，恢复时清空即可，document_outlines / outline_contents
-- 不受影响（不会触发物理删除的 ON DELETE CASCADE）。

ALTER TABLE documents ADD COLUMN IF NOT EXISTS deleted_at TIMESTAMPTZ;

CREATE INDEX IF NOT EXISTS idx_documents_user_deleted
    ON documents (user_id, deleted_at);

COMMENT ON COLUMN documents.deleted_at IS '逻辑删除时间戳；NULL 表示未删除，可通过恢复接口清空';
