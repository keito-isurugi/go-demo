CREATE TABLE IF NOT EXISTS categories
(
    id   SERIAL PRIMARY KEY NOT NULL,
    name VARCHAR(255)       NOT NULL,
    created_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP    NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP    NULL
);

COMMENT ON TABLE categories IS 'カテゴリテーブル';
COMMENT ON COLUMN categories.id IS 'ID';
COMMENT ON COLUMN categories.name IS 'カテゴリ名';
COMMENT ON COLUMN categories.created_at IS '登録日時';
COMMENT ON COLUMN categories.updated_at IS '更新日時';
COMMENT ON COLUMN categories.updated_at IS '削除日時';