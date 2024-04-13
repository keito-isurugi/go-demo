CREATE TABLE IF NOT EXISTS users
(
    id         SERIAL PRIMARY KEY NOT NULL,
    name       VARCHAR(255)       NOT NULL,
    email      VARCHAR(255)       NOT NULL,
    password   VARCHAR(255)       NOT NULL,
    created_at TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP          NOT NULL DEFAULT CURRENT_TIMESTAMP,
    deleted_at TIMESTAMP          NULL
);

COMMENT ON TABLE users IS 'ユーザーテーブル';
COMMENT ON COLUMN users.id IS 'ID';
COMMENT ON COLUMN users.name IS 'ユーザー名';
COMMENT ON COLUMN users.email IS 'メールアドレス';
COMMENT ON COLUMN users.password IS 'パスワード';
COMMENT ON COLUMN users.created_at IS '登録日時';
COMMENT ON COLUMN users.updated_at IS '更新日時';
COMMENT ON COLUMN users.updated_at IS '削除日時';