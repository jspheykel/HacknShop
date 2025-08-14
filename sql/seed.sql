-- 2) Seed data
INSERT INTO categories (name, description) VALUES
('Action','Fast-paced action games'),
('RPG','Role-playing games'),
('Sports','Sports and simulation'),
('Indie','Independent titles');

-- For demo only – store bcrypt hashes in real use.
-- Passwords here are plain text just to unblock development.
INSERT INTO users (username, email, password_hash, is_admin)
VALUES
('admin', 'admin@example.com', '$2a$10$adminhashplaceholder', TRUE),
('heykel', 'heykel@example.com', '$2a$10$userhashplaceholder', FALSE);

INSERT INTO games (title, category_id, description, price_cents, stock, is_active) VALUES
('Elden Ring', 2, 'Open-world action RPG', 899900, 10, TRUE),
('FIFA 24', 3, 'Football simulation', 699900, 15, TRUE),
('Hades', 4, 'Rogue-like dungeon crawler', 299900, 25, TRUE),
('God of War', 1, 'Action adventure', 799900, 8, TRUE),
('Stardew Valley', 4, 'Farming RPG', 149900, 40, TRUE),
('Forza Horizon 5', 3, 'Racing', 749900, 12, TRUE),
('The Witcher 3', 2, 'RPG classic', 399900, 20, TRUE),
('Celeste', 4, 'Indie platformer', 199900, 30, TRUE),
('Sekiro', 1, 'Action adventure', 599900, 9, TRUE),
('NBA 2K24', 3, 'Basketball sim', 659900, 18, TRUE);