TRUNCATE TABLE users RESTART IDENTITY CASCADE;
INSERT INTO users (name, email, password)
VALUES ('山田太郎', 'yamada@example.com', 'pass'),
       ('鈴木花子', 'suzuki@example.com', 'pass'),
       ('渡辺元太', 'watanabe@example.com', 'pass'),
       ('佐藤佳子', 'satou@example.com', 'pass');

TRUNCATE TABLE todos RESTART IDENTITY CASCADE;
INSERT INTO todos (title, done_flag)
VALUES ('テストToDo1', 'false'),
       ('テストToDo2', 'false'),
       ('テストToDo3', 'false'),
       ('テストToDo4', 'false'),
       ('テストToDo5', 'false'),
       ('テストToDo6', 'false'),
       ('テストToDo7', 'false'),
       ('テストToDo8', 'false'),
       ('テストToDo9', 'false'),
       ('テストToDo10', 'false');

TRUNCATE TABLE categories RESTART IDENTITY CASCADE;
INSERT INTO categories (name)
VALUES ('仕事'), ('プライベート'), ('ショッピング'), ('習い事'), ('家事'), ('勉強');

TRUNCATE TABLE todo_categories RESTART IDENTITY CASCADE;
INSERT INTO todo_categories (todo_id, category_id)
VALUES (1, 1), (1, 2), (1, 3), (2, 4), (2, 4), (3, 1);
