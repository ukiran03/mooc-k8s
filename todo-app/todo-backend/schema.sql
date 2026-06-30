-- Force a complete reset every time this script runs
DROP TABLE IF EXISTS tasks;

-- Recreate Table clean
CREATE TABLE tasks (
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    state INTEGER NOT NULL DEFAULT 0 CHECK (state IN (0, 1))
);

-- Always inject the exact same demo dataset
INSERT INTO tasks (title, state) VALUES
('Buy groceries (milk, eggs, bread)', 0),
('Finish reading Chapter 3 of my book', 1),
('Call mom', 1),
('Update resume and LinkedIn profile', 0);
