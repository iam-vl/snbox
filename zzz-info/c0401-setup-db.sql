mysql -u root -p
CREATE DATABASE snbox character set utf8mb4 collate utf8mb4_unicode_ci;
USE snbox;

create table snippets (
    id INTEGER NOT NULL PRIMARY KEY AUTO_INCREMENT,
    title VARCHAR(100) NOT NULL,
    content TEXT NOT NULL,
    created DATETIME NOT NULL,
    expires DATETIME NOT NULL
);

-- Add an index on the created column
CREATE INDEX idx_snippets_created ON snippets(created);

INSERT INTO snippets (title, content, created, expires) VALUES (
    'An old silent pond',
    'An old silent pond...\nA frog jumps into the pond,\nsplash! Silence again.\n\n– Matsuo Bashō', 
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title,content,created,expires) VALUES (
    'Over the wintry forest',
    'Over the wintry\nforest, winds howl in rage\nwith no leaves to blow.\n\n– Natsume Soseki',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 365 DAY)
);

INSERT INTO snippets (title,content,created,expires) VALUES (
    'First autumn morning',
    'First autumn morning\nthe mirror I stare into\nshows my father''s face.\n\n– Murakami Kijo',
    UTC_TIMESTAMP(),
    DATE_ADD(UTC_TIMESTAMP(), INTERVAL 7 DAY)
);

-- CREATE A NEW USER
CREATE USER 'web'@'localhost' IDENTIFIED WITH mysql_native_password BY 'vl#123pass';
GRANT SELECT, INSERT, UPDATE, DELETE ON snbox.* TO 'web'@'localhost';
-- ALTER USER 'web'@'localhost' IDENTIFIED BY 'vl#123pass';

-- mysql -D snbox -u web -p

SHOW COLUMNS FROM snippets;
select id, title, expires from snippets;

DROP TABLE snippets;
-- ERROR 1142 (42000): DROP command denied to user 'web'@'localhost' for table 'snippets'


