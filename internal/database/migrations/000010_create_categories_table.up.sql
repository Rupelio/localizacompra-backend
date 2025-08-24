CREATE TABLE categories (
    id BIGSERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    parent_id BIGINT,
    FOREIGN KEY (parent_id) REFERENCES categories(id)
);
