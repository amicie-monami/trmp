CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    password TEXT NOT NULL
);

CREATE TABLE IF NOT EXISTS writers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL,
    portrait_url TEXT,
    tags TEXT,
    lifespan TEXT,
    country TEXT,
    occupation TEXT,
    is_favorite BOOLEAN DEFAULT FALSE,
    content TEXT
);

CREATE INDEX IF NOT EXISTS idx_writers_name ON writers(name);

-- -- писатели
-- CREATE TABLE IF NOT EXISTS writers (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     name TEXT NOT NULL,
--     portrait_url TEXT,
--     tag TEXT,
--     lifespan TEXT,
--     country TEXT,
--     occupation TEXT,
--     content TEXT,
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP
-- );

-- -- статьи
-- CREATE TABLE IF NOT EXISTS articles (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     title TEXT NOT NULL,
--     tag TEXT,
--     description TEXT,
--     cover_url TEXT,
--     content TEXT,
--     writer_id INTEGER,
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (writer_id) REFERENCES writers (id) ON DELETE SET NULL
-- );

-- -- избранные писатели
-- CREATE TABLE IF NOT EXISTS favorite_writers (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     user_id INTEGER NOT NULL,
--     writer_id INTEGER NOT NULL,
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
--     FOREIGN KEY (writer_id) REFERENCES writers (id) ON DELETE CASCADE,
--     UNIQUE(user_id, writer_id)
-- );

-- -- избранные статьи
-- CREATE TABLE IF NOT EXISTS favorite_articles (
--     id INTEGER PRIMARY KEY AUTOINCREMENT,
--     user_id INTEGER NOT NULL,
--     article_id INTEGER NOT NULL,
--     created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
--     FOREIGN KEY (user_id) REFERENCES users (id) ON DELETE CASCADE,
--     FOREIGN KEY (article_id) REFERENCES articles (id) ON DELETE CASCADE,
--     UNIQUE(user_id, article_id)
-- );

-- CREATE INDEX IF NOT EXISTS idx_favorite_writers_user_id ON favorite_writers(user_id);
-- CREATE INDEX IF NOT EXISTS idx_favorite_writers_writer_id ON favorite_writers(writer_id);
-- CREATE INDEX IF NOT EXISTS idx_favorite_articles_user_id ON favorite_articles(user_id);
-- CREATE INDEX IF NOT EXISTS idx_favorite_articles_article_id ON favorite_articles(article_id);
-- CREATE INDEX IF NOT EXISTS idx_articles_writer_id ON articles(writer_id);