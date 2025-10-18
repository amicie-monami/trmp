INSERT OR IGNORE INTO users (username, email, password_hash) VALUES 
('testuser', 'test@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi'), -- password
('admin', 'admin@example.com', '$2a$10$92IXUNpkjO0rOQ5byMi.Ye4oKoEa3Ro9llC/.og/at2.uheWG/igi');

INSERT OR IGNORE INTO writers (name, portrait_url, lifespan, country, occupation, content) VALUES 
('Лев Толстой', '/images/tolstoy.jpg', '1828-1910', 'Россия', 'Писатель', 'Великий русский писатель...'),
('Фёдор Достоевский', '/images/dostoevsky.jpg', '1821-1881', 'Россия', 'Писатель', 'Известный русский писатель...'),
('Антон Чехов', '/images/chekhov.jpg', '1860-1904', 'Россия', 'Писатель, Драматург', 'Мастер короткого рассказа...');

INSERT OR IGNORE INTO articles (title, description, cover_url, content, writer_id) VALUES 
('Анализ "Войны и мира"', 'Глубокий анализ романа Льва Толстого', '/covers/war-and-peace.jpg', 'Полный текст анализа...', 1),
('Преступление и наказание - разбор', 'Анализ романа Достоевского', '/covers/crime-and-punishment.jpg', 'Разбор основных тем...', 2),
('Чеховские рассказы', 'Особенности коротких рассказов Чехова', '/covers/chekhov-stories.jpg', 'Анализ стиля...', 3);