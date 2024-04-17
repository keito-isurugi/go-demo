TRUNCATE TABLE users RESTART IDENTITY CASCADE;
INSERT INTO users (name, email, password)
VALUES ('山田太郎', 'yamada@example.com', 'pass'),
       ('鈴木花子', 'suzuki@example.com', 'pass'),
       ('渡辺元太', 'watanabe@example.com', 'pass'),
       ('佐藤佳子', 'satou@example.com', 'pass');

TRUNCATE TABLE todos RESTART IDENTITY CASCADE;
INSERT INTO todos (user_id, title, done_flag)
VALUES (1, 'テストToDo1', 'false'),
       (1, 'テストToDo2', 'false'),
       (1, 'テストToDo3', 'false'),
       (1, 'テストToDo4', 'false'),
       (2, 'テストToDo5', 'false'),
       (2, 'テストToDo6', 'false'),
       (2, 'テストToDo7', 'false'),
       (3, 'テストToDo8', 'false'),
       (3, 'テストToDo9', 'false'),
       (3, 'テストToDo10', 'false');

TRUNCATE TABLE categories RESTART IDENTITY CASCADE;
INSERT INTO categories (name)
VALUES ('仕事'), ('プライベート'), ('ショッピング'), ('習い事'), ('家事'), ('勉強');

TRUNCATE TABLE todo_categories RESTART IDENTITY CASCADE;
INSERT INTO todo_categories (todo_id, category_id)
VALUES (1, 1), (1, 2), (1, 3), (2, 4), (2, 4), (3, 1);
