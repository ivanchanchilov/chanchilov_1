BEGIN;
CREATE TABLE categories (
    id INTEGER PRIMARY KEY,
    name TEXT NOT NULL
);
INSERT INTO categories (id, name) VALUES
    (1, 'Другое'),
    (2, 'Еда'),
    (3, 'Транспорт'),
    (4, 'Жильё'),
    (5, 'Развлечения');
CREATE TABLE expenses_new (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    amount REAL NOT NULL,
    description TEXT NOT NULL,
    date TEXT NOT NULL,
    category_id INTEGER NOT NULL DEFAULT 1 REFERENCES categories(id)
);
INSERT INTO expenses_new (id, amount, description, date, category_id)
SELECT id, amount, description, date, 1 FROM expenses;
DROP TABLE expenses;
ALTER TABLE expenses_new RENAME TO expenses;
PRAGMA user_version = 2;
COMMIT;
