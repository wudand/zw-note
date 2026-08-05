-- Migration: 003_create_documents_table
-- Direction: UP

CREATE TABLE IF NOT EXISTS documents (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT       NOT NULL,
    title       VARCHAR(50)  NOT NULL,
    description VARCHAR(200),
    author      VARCHAR(20)  NOT NULL,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_documents_user
        FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE RESTRICT
);

CREATE INDEX IF NOT EXISTS idx_documents_user_updated
    ON documents (user_id, updated_at DESC);

DROP TRIGGER IF EXISTS trg_documents_updated_at ON documents;
CREATE TRIGGER trg_documents_updated_at
    BEFORE UPDATE ON documents
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE documents IS '知识笔记文档元信息';
COMMENT ON COLUMN documents.user_id IS '所属用户（所有权根）';
COMMENT ON COLUMN documents.author IS '展示作者，可与 users.display_name 不同';
