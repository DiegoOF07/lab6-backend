CREATE TABLE IF NOT EXISTS series(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    title TEXT NOT NULL UNIQUE,
    status TEXT NOT NULL,
    episodes INTEGER NOT NULL,
    last_episode INTEGER CHECK(last_episode<=episodes),
    ranking INTEGER NOT NULL
);
