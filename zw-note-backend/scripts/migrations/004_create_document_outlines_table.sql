-- Migration: 004_create_document_outlines_table
-- Direction: UP

CREATE TABLE IF NOT EXISTS document_outlines (
    id          BIGSERIAL PRIMARY KEY,
    document_id BIGINT       NOT NULL,
    parent_id   BIGINT,
    title       VARCHAR(100) NOT NULL,
    sort_order  INT          NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ  NOT NULL DEFAULT NOW(),
    CONSTRAINT fk_outlines_document
        FOREIGN KEY (document_id) REFERENCES documents (id) ON DELETE CASCADE,
    CONSTRAINT fk_outlines_parent
        FOREIGN KEY (parent_id) REFERENCES document_outlines (id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_outlines_document_parent_sort
    ON document_outlines (document_id, parent_id, sort_order);

CREATE INDEX IF NOT EXISTS idx_outlines_document_id
    ON document_outlines (document_id);

DROP TRIGGER IF EXISTS trg_document_outlines_updated_at ON document_outlines;
CREATE TRIGGER trg_document_outlines_updated_at
    BEFORE UPDATE ON document_outlines
    FOR EACH ROW EXECUTE FUNCTION set_updated_at();

COMMENT ON TABLE document_outlines IS '文档目录树节点';
COMMENT ON COLUMN document_outlines.parent_id IS '父节点 ID；NULL 表示顶级';
COMMENT ON COLUMN document_outlines.sort_order IS '同级排序，越小越靠前';
