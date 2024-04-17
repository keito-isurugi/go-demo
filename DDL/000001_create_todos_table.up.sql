CREATE TABLE IF NOT EXISTS todos
(
    id         SERIAL PRIMARY KEY NOT NULL,
    user_id    INT                 NOT NULL,
    title      VARCHAR(255)        NOT NULL,
    done_flag  BOOLEAN             NOT NULL DEFAULT false,
    created_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP           NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP           NULL
);

COMMENT ON TABLE todos IS 'Todosテーブル';
COMMENT ON COLUMN todos.id IS 'ID';
COMMENT ON COLUMN todos.user_id IS 'ユーザーID';
COMMENT ON COLUMN todos.title IS 'タイトル';
COMMENT ON COLUMN todos.done_flag IS '完了フラグ';
COMMENT ON COLUMN todos.created_at IS '登録日時';
COMMENT ON COLUMN todos.updated_at IS '更新日時';
COMMENT ON COLUMN todos.updated_at IS '削除日時';