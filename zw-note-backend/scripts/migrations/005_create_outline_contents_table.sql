-- Migration: 005_create_outline_contents_table
-- Direction: UP

CREATE TABLE IF NOT EXISTS outline_contents (
    outline_id BIGINT PRIMARY KEY,
    content    TEXT        NOT NULL DEFAULT '',
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_contents_outline
        FOREIGN KEY (outline_id) REFERENCES document_outlines (id) ON DELETE CASCADE
);

DROP TRIGGER IF EXISTS trg_outline_contents_updated_at ON outline_contents;
CREATE TRIGGER trg_outline_contents_updated_at
    BEFORE UPDATE ON outline_contents
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE outline_contents IS '目录节点 Markdown 内容（与 outline 1:1）';
COMMENT ON COLUMN outline_contents.outline_id IS '主键即目录节点 ID，无独立 surrogate id';
