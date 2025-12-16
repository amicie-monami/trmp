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


CREATE TABLE IF NOT EXISTS articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    cover_url TEXT,
    title TEXT NOT NULL,
    tags TEXT,
    description TEXT,
    is_favorite BOOLEAN DEFAULT FALSE,
    content TEXT
);

-- избранное
CREATE TABLE IF NOT EXISTS favorite_writers (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    writer_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (writer_id) REFERENCES writers(id) ON DELETE CASCADE,
    UNIQUE(user_id, writer_id)  
);

CREATE TABLE IF NOT EXISTS favorite_articles (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    user_id INTEGER NOT NULL,
    article_id INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY (user_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY (article_id) REFERENCES articles(id) ON DELETE CASCADE,
    UNIQUE(user_id, article_id)  
);

CREATE INDEX IF NOT EXISTS idx_favorite_writers_user ON favorite_writers(user_id);
CREATE INDEX IF NOT EXISTS idx_favorite_writers_writer ON favorite_writers(writer_id);
CREATE INDEX IF NOT EXISTS idx_favorite_articles_user ON favorite_articles(user_id);
CREATE INDEX IF NOT EXISTS idx_favorite_articles_article ON favorite_articles(article_id);
