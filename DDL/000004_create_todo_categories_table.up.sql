CREATE TABLE IF NOT EXISTS todo_categories
(
    id          SERIAL PRIMARY KEY NOT NULL,
    todo_id     INT                 NOT NULL,
    category_id INT                 NOT NULL,
    created_at TIMESTAMP            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP            NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP            NULL
);

COMMENT ON TABLE todo_categories IS 'Todo, カテゴリ中間テーブル';
COMMENT ON COLUMN todo_categories.id IS 'ID';
COMMENT ON COLUMN todo_categories.todo_id IS 'TodoID';
COMMENT ON COLUMN todo_categories.category_id IS 'カテゴリID';
COMMENT ON COLUMN todo_categories.created_at IS '登録日時';
COMMENT ON COLUMN todo_categories.updated_at IS '更新日時';
COMMENT ON COLUMN todo_categories.updated_at IS '削除日時';