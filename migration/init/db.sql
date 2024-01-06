CREATE TABLE url(
    id SERIAL PRIMARY KEY,
    original_url VARCHAR(2048) NOT NULL,
    short_url VARCHAR(20) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

INSERT INTO url(original_url, short_url) VALUES
    ('https://rezapu.dev', 'abc'),
    ('https://example.com', 'example');